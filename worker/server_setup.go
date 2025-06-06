package worker

import (
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"mt-hosting-manager/worker/server_setup"
	"time"
)

func (w *Worker) ServerSetup(ctx *JobContext) error {
	job := ctx.job
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
			job.Step = 2
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

		restore_zip_filename := "restore.zip"

		// try to remove file if it exists
		client.DeleteFile(restore_zip_filename)

		err = client.SetMaintenanceMode(true)
		if err != nil {
			return fmt.Errorf("set maintenance mode error: %v", err)
		}

		uw := client.UploadStream(restore_zip_filename, func(written_bytes int64) {
			job.ProgressPercent = float64(written_bytes) / float64(backup.Size) * 100
			job.Message = fmt.Sprintf("Transferred bytes: %d / %d (%.2f %%)", written_bytes, backup.Size, job.ProgressPercent)
			ctx.w.repos.JobRepo.UpdateWithTx(ctx.tx, job)
		})
		err = w.core.StreamBackup(backup, uw)
		if err != nil {
			return fmt.Errorf("stream error: %v", err)
		}
		uw.Close()

		err = client.UnzipFile(restore_zip_filename)
		if err != nil {
			return fmt.Errorf("unzip error: %v", err)
		}

		err = client.DeleteFile(restore_zip_filename)
		if err != nil {
			return fmt.Errorf("restore-file cleanup error: %v", err)
		}

		err = client.SetMaintenanceMode(false)
		if err != nil {
			return fmt.Errorf("set maintenance mode error: %v", err)
		}

		job.Message = "Restore complete, exiting maintenance mode"
		job.NextRun = time.Now().Add(5 * time.Second).Unix()
		job.Step = 2

	case 2:
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
