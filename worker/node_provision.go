package worker

import (
	"fmt"
	"mt-hosting-manager/api/hetzner_cloud"
	"mt-hosting-manager/api/hetzner_dns"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/provision"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func (w *Worker) NodeProvision(job *types.Job) error {
	node, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return fmt.Errorf("node not found: %s", *job.UserNodeID)
	}

	nodetype, err := w.repos.NodeTypeRepo.GetByID(node.NodeTypeID)
	if err != nil {
		return err
	}
	if nodetype == nil {
		return fmt.Errorf("nodetype not found: %s", node.NodeTypeID)
	}

	switch job.Step {
	case 0:
		// hetzner setup
		w.core.AddAuditLog(&types.AuditLog{
			Type:       types.AuditLogNodeProvisioningStarted,
			UserID:     node.UserID,
			UserNodeID: &node.ID,
		})

		node.State = types.UserNodeStateProvisioning
		err = w.repos.UserNodeRepo.Update(node)
		if err != nil {
			return err
		}

		if node.ExternalID == "" {
			// create new instance
			logrus.WithFields(logrus.Fields{
				"node_id": node.ID,
			}).Info("Creating new hetzner server instance")

			job.Message = "creating new server"
			job.ProgressPercent = 10

			csr := &hetzner_cloud.CreateServerRequest{
				Image: "ubuntu-22.04",
				Labels: map[string]string{
					"node_id": node.ID,
					"stage":   w.cfg.Stage,
				},
				Location: node.Location,
				Name:     node.Name,
				PublicNet: &hetzner_cloud.PublicNet{
					EnableIPv4: true,
					EnableIPv6: true,
				},
				ServerType:       nodetype.ServerType,
				SSHKeys:          []string{"minetest@keymaster", "thomas@keymaster"},
				StartAfterCreate: true,
			}
			resp, err := w.hcc.CreateServer(csr)
			if err != nil {
				return fmt.Errorf("create server failed: %v", err)
			}
			logrus.WithFields(logrus.Fields{
				"hetzner_id": resp.Server.ID,
				"status":     resp.Server.Status,
			}).Info("Server created")

			node.ExternalID = fmt.Sprintf("%d", resp.Server.ID)
			node.IPv4 = resp.Server.PublicNet.IPv4.IP

			// replace trailing /64 with 1
			ipv6_parts := strings.Split(resp.Server.PublicNet.IPv6.IP, "/")
			node.IPv6 = fmt.Sprintf("%s1", ipv6_parts[0])
			err = w.repos.UserNodeRepo.Update(node)
			if err != nil {
				return fmt.Errorf("external-id update failed: %v", err)
			}
		}

		job.Step = 1
	case 1:
		if node.ExternalIPv4DNSID == "" {
			job.Message = "creating dns A-record"
			job.ProgressPercent = 20

			record_name := fmt.Sprintf("%s%s", node.Name, w.cfg.DNSRecordSuffix)
			record, err := w.CreateDNSRecord(hetzner_dns.RecordA, record_name, node.IPv4)
			if err != nil {
				return fmt.Errorf("could not create A-record: %v", err)
			}
			node.ExternalIPv4DNSID = record.ID
			err = w.repos.UserNodeRepo.Update(node)
			if err != nil {
				return fmt.Errorf("ipv4 record update failed: %v", err)
			}
		}
		job.Step = 2

	case 2:
		if node.ExternalIPv6DNSID == "" {
			job.Message = "creating dns AAAA-record"
			job.ProgressPercent = 30

			record_name := fmt.Sprintf("%s%s", node.Name, w.cfg.DNSRecordSuffix)
			record, err := w.CreateDNSRecord(hetzner_dns.RecordAAAA, record_name, node.IPv6)
			if err != nil {
				return fmt.Errorf("could not create AAAA-record: %v", err)
			}
			node.ExternalIPv6DNSID = record.ID
			err = w.repos.UserNodeRepo.Update(node)
			if err != nil {
				return fmt.Errorf("ipv6 record update failed: %v", err)
			}
		}
		job.Step = 3

	case 3:
		job.Message = "waiting for node-startup"
		job.ProgressPercent = 40
		job.NextRun = time.Now().Add(time.Second * 20).Unix()
		job.Step = 4

	case 4:
		client, err := core.TrySSHConnection(node)
		if err != nil {
			return fmt.Errorf("ssh try error: %v", err)
		}
		client.Close()
		job.Message = "Executing provisioning script"
		job.ProgressPercent = 60
		job.Step = 5

	case 5:
		logrus.WithFields(logrus.Fields{
			"external_id": node.ExternalID,
			"ipv4":        node.IPv4,
			"ipv6":        node.IPv6,
		}).Info("Executing provisioning")

		client, err := core.TrySSHConnection(node)
		if err != nil {
			return fmt.Errorf("ssh try error: %v", err)
		}
		err = provision.Provision(client, w.cfg, node.UserID)
		if err != nil {
			return fmt.Errorf("provision error: %v", err)
		}
		err = client.Close()
		if err != nil {
			return fmt.Errorf("client close error: %v", err)
		}

		node.State = types.UserNodeStateRunning
		err = w.repos.UserNodeRepo.Update(node)
		if err != nil {
			return fmt.Errorf("usernode update error: %v", err)
		}

		w.core.AddAuditLog(&types.AuditLog{
			Type:       types.AuditLogNodeProvisioningFinished,
			UserID:     node.UserID,
			UserNodeID: &node.ID,
		})

		job.State = types.JobStateDoneSuccess
		job.ProgressPercent = 100
		job.Message = "done"
	}

	return nil
}
