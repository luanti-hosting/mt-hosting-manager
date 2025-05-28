package worker

import (
	"fmt"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/server_setup"
	"time"
)

func (w *Worker) serverPrepareSetup(node *types.UserNode, server *types.MinetestServer) error {

	server.State = types.MinetestServerStateProvisioning
	err := w.repos.MinetestServerRepo.Update(server)
	if err != nil {
		return fmt.Errorf("server entity update error: %v", err)
	}

	record_name := fmt.Sprintf("%s%s", server.DNSName, w.cfg.DNSRecordSuffix)
	record_value := fmt.Sprintf("%s%s", node.Name, w.cfg.DNSRecordSuffix)
	if server.ExternalCNAMEDNSID == "" {
		// create new record
		record, err := w.CreateDNSRecord(hetzner_dns.RecordCNAME, record_name, record_value)
		if err != nil {
			return fmt.Errorf("could not create CNAME record: %v", err)
		}
		server.ExternalCNAMEDNSID = record.ID

	} else {
		// check if record matches config
		record, err := w.hdc.GetRecord(server.ExternalCNAMEDNSID)
		if err != nil {
			return fmt.Errorf("could not fetch current cname with id: '%s': %v", server.ExternalCNAMEDNSID, err)
		}
		if record.Record.Name != record_name || record.Record.Value != record_value {
			// values changed, remove and recreate
			err = w.hdc.DeleteRecord(server.ExternalCNAMEDNSID)
			if err != nil {
				return fmt.Errorf("could not remove record with id '%s', %v", server.ExternalCNAMEDNSID, err)
			}
			created_record, err := w.CreateDNSRecord(hetzner_dns.RecordCNAME, record_name, record_value)
			if err != nil {
				return fmt.Errorf("could not re-create CNAME record: %v", err)
			}
			server.ExternalCNAMEDNSID = created_record.ID
		}
	}

	// save external dns id
	err = w.repos.MinetestServerRepo.Update(server)
	if err != nil {
		return fmt.Errorf("mid-setup update failed: %v", err)
	}

	// dns propagation time (LE has issues with really _fresh_ records)
	time.Sleep(20 * time.Second)

	return nil
}

// removes a server instance and optionally removes all the containing data
func (w *Worker) removeServer(node *types.UserNode, server *types.MinetestServer, cleanup_data bool) error {
	server.State = types.MinetestServerStateRemoving
	err := w.repos.MinetestServerRepo.Update(server)
	if err != nil {
		return fmt.Errorf("server entity update error: %v", err)
	}

	if server.ExternalCNAMEDNSID != "" {
		err = w.hdc.DeleteRecord(server.ExternalCNAMEDNSID)
		if err != nil {
			return fmt.Errorf("could not remove cname (id: %s) of server %s: %v", server.ExternalCNAMEDNSID, server.DNSName, err)
		}
		server.ExternalCNAMEDNSID = ""
		err = w.repos.MinetestServerRepo.Update(server)
		if err != nil {
			return fmt.Errorf("could not update server entry '%s': %v", server.ID, err)
		}
	}

	if cleanup_data {
		client, err := core.TrySSHConnection(node)
		if err != nil {
			return err
		}

		// remove potentially running services
		_, _, err = core.SSHExecute(client, fmt.Sprintf("docker rm -f %s || true", server_setup.GetEngineName(server)))
		if err != nil {
			return fmt.Errorf("could not stop running service: %v", err)
		}

		basedir := server_setup.GetBaseDir(server)
		_, _, err = core.SSHExecute(client, fmt.Sprintf("cd %s && docker compose down -v", basedir))
		if err != nil {
			return fmt.Errorf("could not run docker compose down: %v", err)
		}

		_, _, err = core.SSHExecute(client, fmt.Sprintf("rm -rf %s", basedir))
		if err != nil {
			return fmt.Errorf("could not run remove data-dir '%s': %v", basedir, err)
		}
	}

	server.State = types.MinetestServerStateDecommissioned
	err = w.repos.MinetestServerRepo.Update(server)
	if err != nil {
		return fmt.Errorf("server entity update error (decommissioned): %v", err)
	}

	w.core.AddAuditLog(&types.AuditLog{
		Type:             types.AuditLogServerRemoved,
		UserID:           node.UserID,
		UserNodeID:       &node.ID,
		MinetestServerID: &server.ID,
	})

	return nil
}
