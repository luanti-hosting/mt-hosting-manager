package worker

import (
	"fmt"
	"mt-hosting-manager/api/mtui"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/server_setup"
	"time"
)

func (w *Worker) ServerSetup(job *types.Job) error {
	node, server, err := w.GetJobContext(job)
	if err != nil {
		return err
	}

	switch job.Step {
	case 0:
		w.core.AddAuditLog(&types.AuditLog{
			Type:             types.AuditLogServerSetupStarted,
			UserID:           node.UserID,
			UserNodeID:       &node.ID,
			MinetestServerID: &server.ID,
		})

		err = w.serverPrepareSetup(node, server)
		if err != nil {
			return err
		}

		client, err := core.TrySSHConnection(node)
		if err != nil {
			return err
		}
		defer client.Close()

		err = server_setup.Setup(client, w.cfg, node, server)
		if err != nil {
			return err
		}

		if job.BackupID == nil {
			// skip restore steps
			job.Step = 3
			return nil
		} else {
			// restore after the tls connection can be established
			job.Step = 1
			job.Message = "Restore pending"
			job.NextRun = time.Now().Add(60 * time.Second).Unix()
			return nil
		}

	case 1:
		// trigger restore

		client, err := w.core.GetMTUIClient(server)
		if err != nil {
			return fmt.Errorf("get client error: %v", err)
		}

		backup, err := w.repos.BackupRepo.GetByID(*job.BackupID)
		if err != nil {
			return fmt.Errorf("get backup error: %v", err)
		}
		if backup == nil {
			return fmt.Errorf("backup not found: '%s'", *job.BackupID)
		}

		// go into maintenance mode
		err = client.SetMaintenanceMode(true)
		if err != nil {
			return fmt.Errorf("could not enable maintenance mode: %v", err)
		}

		// create restore job
		info, err := client.CreateBackupRestoreJob(&mtui.CreateBackupRestoreJob{
			Type:      mtui.RestoreJob,
			Endpoint:  w.cfg.S3Endpoint,
			KeyID:     w.cfg.S3KeyID,
			AccessKey: w.cfg.S3AccessKey,
			Bucket:    w.cfg.S3Bucket,
			FileKey:   backup.Passphrase,
			Filename:  fmt.Sprintf("%s.zip", backup.ID),
		})
		if err != nil {
			return fmt.Errorf("create restore job error: %v", err)
		}

		job.Message = info.Message
		job.Step = 2
		job.NextRun = time.Now().Add(2 * time.Second).Unix()

	case 2:
		// get restore status

		client, err := w.core.GetMTUIClient(server)
		if err != nil {
			return fmt.Errorf("get client error: %v", err)
		}

		// check restore job
		info, err := client.GetBackupRestoreJobInfo()
		if err != nil {
			return fmt.Errorf("get restore job error: %v", err)
		}

		switch info.State {
		case mtui.BackupRestoreJobRunning:
			// still running
			job.Message = info.Message
			job.NextRun = time.Now().Add(2 * time.Second).Unix()

		case mtui.BackupRestoreJobSuccess:
			// all done
			job.Message = info.Message
			job.Step = 3

		case mtui.BackupRestoreJobFailure:
			// restore failed
			job.Message = info.Message
			job.State = types.JobStateDoneFailure
		}

	case 3:
		// mark running

		server.State = types.MinetestServerStateRunning
		err = w.repos.MinetestServerRepo.Update(server)
		if err != nil {
			return fmt.Errorf("server entity update error: %v", err)
		}

		w.core.AddAuditLog(&types.AuditLog{
			Type:             types.AuditLogServerSetupFinished,
			UserID:           node.UserID,
			UserNodeID:       &node.ID,
			MinetestServerID: &server.ID,
		})

		job.State = types.JobStateDoneSuccess
	}

	return nil
}
