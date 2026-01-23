package worker

import (
	"context"
	"errors"
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/sirupsen/logrus"
)

func (w *Worker) NodeDestroy(job *types.Job) error {
	node, err := w.repos.UserNodeRepo.GetByID(*job.UserNodeID)
	if err != nil {
		return err
	}
	if node == nil {
		return errors.New("node not found")
	}

	nt, err := w.repos.NodeTypeRepo.GetByID(node.NodeTypeID)
	if err != nil {
		return err
	}
	if nt == nil {
		return errors.New("node-type not found")
	}

	user, err := w.repos.UserRepo.GetByID(node.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	servers, err := w.repos.MinetestServerRepo.Search(&types.MinetestServerSearch{NodeID: &node.ID})
	if err != nil {
		return fmt.Errorf("could not fetch servers: %v", err)
	}

	jd := &types.RemoveNodeJobData{}
	if job.HasData() {
		err = job.GetData(jd)
		if err != nil {
			return fmt.Errorf("could get job data: %v", err)
		}
	}

	if !jd.CreateBackups {
		// skip backup step
		job.Step = 2
	}

	switch job.Step {
	case 0: // create backup job(s)
		for _, server := range servers {
			if server.State != types.MinetestServerStateRunning {
				// server is already removed, no need for backup
				continue
			}

			now := time.Now()
			backup := &types.Backup{
				ID:               uuid.NewString(),
				State:            types.BackupStateCreated,
				Passphrase:       core.RandStringRunes(64),
				UserID:           node.UserID,
				MinetestServerID: server.ID,
				Created:          now.Unix(),
				Expires:          now.Add(time.Hour * 24 * 365).Unix(),
			}
			err = w.repos.BackupRepo.Insert(backup)
			if err != nil {
				return fmt.Errorf("backup insert error: %v", err)
			}

			backup_job := types.BackupServerJob(node, server, backup)
			err = w.repos.JobRepo.Insert(backup_job)
			if err != nil {
				return fmt.Errorf("backup-job insert error: %v", err)
			}

			jd.BackupJobIDs = append(jd.BackupJobIDs, backup_job.ID)
		}

		err = job.SetData(jd)
		if err != nil {
			return fmt.Errorf("could not set job data: %v", err)
		}
		job.Step = 1
		job.Message = fmt.Sprintf("waiting for %d backup(s) to finish", len(jd.BackupJobIDs))
		job.NextRun = time.Now().Add(5 * time.Second).Unix()

	case 1: // wait for backup jobs to finish

		done_count := 0
		for _, backup_job_id := range jd.BackupJobIDs {
			backup_job, err := w.repos.JobRepo.GetByID(backup_job_id)
			if err != nil {
				return fmt.Errorf("could not get backup job %s: %v", backup_job_id, err)
			}
			switch backup_job.State {
			case types.JobStateDoneFailure:
				return fmt.Errorf("backup-job failed: %s", backup_job_id)
			case types.JobStateDoneSuccess:
				done_count++
			}
		}

		if done_count == len(jd.BackupJobIDs) {
			job.Message = "Backups done"
			job.Step = 2
		} else {
			job.Message = fmt.Sprintf("Backups: %d/%d done", done_count, len(jd.BackupJobIDs))
		}

		job.NextRun = time.Now().Add(5 * time.Second).Unix()

	case 2: // remove everything
		for _, server := range servers {
			if server.State != types.MinetestServerStateRunning {
				// server is already removed
				continue
			}

			err = w.removeServer(node, server, false) // no need to remove data
			if err != nil {
				return fmt.Errorf("error removing server '%s': %v", server.ID, err)
			}
		}

		if node.ExternalID != "" {
			server_id, err := strconv.ParseInt(node.ExternalID, 10, 32)
			if err != nil {
				return fmt.Errorf("could not parse server id: '%s'", node.ExternalID)
			}
			_, _, err = w.hc.Server.DeleteWithResult(context.Background(), &hcloud.Server{ID: server_id})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"ExternalID": node.ExternalID,
				}).Warn("Server instance not found, not deleting anything")
			}
			node.ExternalID = ""
			err = w.repos.UserNodeRepo.Update(node)
			if err != nil {
				return fmt.Errorf("could not update usernode: %v", err)
			}
		}

		if node.ExternalIPv4DNSID != "" {
			_, _, err = w.hc.Zone.DeleteRRSet(context.Background(), &hcloud.ZoneRRSet{
				Zone: &hcloud.Zone{Name: w.cfg.HetznerZoneName},
				ID:   node.ExternalIPv4DNSID,
			})
			if err != nil {
				return fmt.Errorf("could not remove A-record: %v", err)
			}
			node.ExternalIPv4DNSID = ""
			err = w.repos.UserNodeRepo.Update(node)
			if err != nil {
				return fmt.Errorf("could not update usernode: %v", err)
			}
		}

		if node.ExternalIPv6DNSID != "" {
			_, _, err = w.hc.Zone.DeleteRRSet(context.Background(), &hcloud.ZoneRRSet{
				Zone: &hcloud.Zone{Name: w.cfg.HetznerZoneName},
				ID:   node.ExternalIPv6DNSID,
			})
			if err != nil {
				return fmt.Errorf("could not remove AAAA-record: %v", err)
			}
			node.ExternalIPv6DNSID = ""
			err = w.repos.UserNodeRepo.Update(node)
			if err != nil {
				return fmt.Errorf("could not update usernode: %v", err)
			}
		}

		w.core.AddAuditLog(&types.AuditLog{
			Type:       types.AuditLogNodeRemoved,
			UserID:     node.UserID,
			UserNodeID: &node.ID,
		})

		notify.Send(&notify.NtfyNotification{
			Title:    fmt.Sprintf("Node removed by %s (Type: %s)", user.Name, nt.Name),
			Message:  fmt.Sprintf("User: %s, Node-type %s, Name: %s", user.Name, nt.Name, node.Name),
			Priority: 3,
			Tags:     []string{"computer", "x"},
		}, true)

		node.State = types.UserNodeStateDecommissioned
		job.State = types.JobStateDoneSuccess

		err = w.repos.UserNodeRepo.Update(node)
		if err != nil {
			return fmt.Errorf("could not update usernode before finish: %v", err)
		}

		job.State = types.JobStateDoneSuccess
		job.ProgressPercent = 100
		job.Message = "done"
	}

	return nil
}
