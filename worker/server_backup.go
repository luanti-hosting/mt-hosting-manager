package worker

import (
	"fmt"
	"mt-hosting-manager/api/mtui"
	"mt-hosting-manager/types"
	"time"
)

func (w *Worker) ServerBackup(job *types.Job) error {
	server, err := w.repos.MinetestServerRepo.GetByID(*job.MinetestServerID)
	if err != nil {
		return fmt.Errorf("get server error: %v", err)
	}
	if server == nil {
		return fmt.Errorf("server not found")
	}

	backup, err := w.repos.BackupRepo.GetByID(*job.BackupID)
	if err != nil {
		return fmt.Errorf("get backup error: %v", err)
	}
	if backup == nil {
		return fmt.Errorf("backup not found: '%s'", *job.BackupID)
	}

	client, err := w.core.GetMTUIClient(server)
	if err != nil {
		return fmt.Errorf("get client error: %v", err)
	}

	switch job.Step {
	case 0:
		// trigger backup
		info, err := client.CreateBackupRestoreJob(&mtui.CreateBackupRestoreJob{
			Type:      mtui.BackupJob,
			Endpoint:  w.cfg.S3Endpoint,
			KeyID:     w.cfg.S3KeyID,
			AccessKey: w.cfg.S3AccessKey,
			Bucket:    w.cfg.S3Bucket,
			FileKey:   backup.Passphrase,
			Filename:  fmt.Sprintf("%s.zip", backup.ID),
		})
		if err != nil {
			return fmt.Errorf("create backup job error: %v", err)
		}

		w.core.AddAuditLog(&types.AuditLog{
			Type:             types.AuditLogServerBackupStarted,
			UserID:           backup.UserID,
			UserNodeID:       job.UserNodeID,
			MinetestServerID: job.MinetestServerID,
			BackupID:         job.BackupID,
		})

		job.Message = info.Message
		job.Step = 1
		job.NextRun = time.Now().Add(2 * time.Second).Unix()
	case 1:
		// check backup
		info, err := client.GetBackupRestoreJobInfo()
		if err != nil {
			return fmt.Errorf("get backup job error: %v", err)
		}

		switch info.State {
		case mtui.BackupRestoreJobRunning:
			// still running
			job.Message = info.Message
			job.ProgressPercent = info.ProgressPercent
			job.NextRun = time.Now().Add(2 * time.Second).Unix()

		case mtui.BackupRestoreJobSuccess:
			// get size from storage
			size, err := w.core.GetBackupSize(backup)
			if err != nil {
				job.State = types.JobStateDoneFailure
				job.Message = fmt.Sprintf("backup-file stat failed: %v", err)
				backup.State = types.BackupStateError
			} else {
				// everything checks out
				job.State = types.JobStateDoneSuccess
				job.Message = info.Message
				backup.State = types.BackupStateComplete
				backup.Size = size
			}

			err = w.repos.BackupRepo.Update(backup)
			if err != nil {
				return fmt.Errorf("error in backup update: %v", err)
			}

			w.core.AddAuditLog(&types.AuditLog{
				Type:             types.AuditLogServerBackupFinished,
				UserID:           backup.UserID,
				UserNodeID:       job.UserNodeID,
				MinetestServerID: job.MinetestServerID,
				BackupID:         job.BackupID,
			})

		case mtui.BackupRestoreJobFailure:
			// backup failed
			job.Message = info.Message
			job.State = types.JobStateDoneFailure
		}
	}

	return nil
}
