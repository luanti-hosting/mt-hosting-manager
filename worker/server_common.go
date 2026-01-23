package worker

import (
	"context"
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/server_setup"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
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
		record, _, err := w.hc.Zone.CreateRRSet(context.Background(), &hcloud.Zone{Name: w.cfg.HetznerZoneName}, hcloud.ZoneRRSetCreateOpts{
			Name: record_name,
			Type: hcloud.ZoneRRSetTypeCNAME,
			TTL:  hcloud.Ptr(300),
			Records: []hcloud.ZoneRRSetRecord{
				{Value: record_value},
			},
		})
		if err != nil {
			return fmt.Errorf("could not create CNAME record: %v", err)
		}
		server.ExternalCNAMEDNSID = record.RRSet.ID

	} else {
		// check if record matches config
		rrset, _, err := w.hc.Zone.GetRRSetByID(context.Background(), &hcloud.Zone{Name: w.cfg.HetznerZoneName}, server.ExternalCNAMEDNSID)
		if err != nil {
			return fmt.Errorf("could not fetch current cname with id: '%s': %v", server.ExternalCNAMEDNSID, err)
		}
		if len(rrset.Records) != 1 {
			return fmt.Errorf("invalid rrset record length: %d", len(rrset.Records))
		}
		record := rrset.Records[0]

		if rrset.Name != record_name || record.Value != record_value {
			// values changed, remove and recreate
			_, _, err = w.hc.Zone.DeleteRRSet(context.Background(), &hcloud.ZoneRRSet{
				Zone: &hcloud.Zone{Name: w.cfg.HetznerZoneName},
				ID:   server.ExternalCNAMEDNSID,
			})
			if err != nil {
				return fmt.Errorf("could not remove record with id '%s', %v", server.ExternalCNAMEDNSID, err)
			}
			created_record, err := w.CreateOrUpdateDNSRecord(hcloud.ZoneRRSetTypeCNAME, record_name, record_value)
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
		_, _, err = w.hc.Zone.DeleteRRSet(context.Background(), &hcloud.ZoneRRSet{
			Zone: &hcloud.Zone{Name: w.cfg.HetznerZoneName},
			ID:   server.ExternalCNAMEDNSID,
		})
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

		// remove compose services
		basedir := server_setup.GetBaseDir(server)
		_, _, err = core.SSHExecute(client, fmt.Sprintf("cd %s && docker compose down -v", basedir))
		if err != nil {
			return fmt.Errorf("could not run docker compose down: %v", err)
		}

		// cleanup networks
		_, _, err = core.SSHExecute(client, "docker network prune -f")
		if err != nil {
			return fmt.Errorf("could not cleanup networks: %v", err)
		}

		// remove data
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
