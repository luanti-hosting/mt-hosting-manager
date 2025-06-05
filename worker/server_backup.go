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
		info, err := client.CreateBackupJob(&mtui.CreateBackupJob{
			ID:       backup.ID,
			Type:     mtui.BackupJobTypeWEBDAV,
			URL:      w.cfg.StorageURL,
			Username: w.cfg.StorageUsername,
			Password: w.cfg.StoragePassword,
			Filename: fmt.Sprintf("%s.zip", backup.ID),
			Key:      backup.Passphrase,
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
		job.Data = []byte(info.ID)
		job.Step = 1
		job.NextRun = time.Now().Add(5 * time.Second).Unix()
	case 1:
		// check backup
		info, err := client.GetBackupJobInfo(string(job.Data))
		if err != nil {
			return fmt.Errorf("get backup job error: %v", err)
		}

		switch info.Status {
		case mtui.BackupJobRunning:
			// still running
			job.Message = info.Message
			job.NextRun = time.Now().Add(5 * time.Second).Unix()

		case mtui.BackupJobSuccess:
			// all done, delay a bit before checking final backup file
			time.Sleep(10 * time.Second)
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

		case mtui.BackupJobFailure:
			// backup failed
			job.Message = info.Message
			job.State = types.JobStateDoneFailure
		}
	}

	return nil
}
