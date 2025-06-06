package worker

import (
	"fmt"
	"mt-hosting-manager/types"
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

	r, err := client.DownloadZip("/")
	if err != nil {
		return fmt.Errorf("download zip error: %v", err)
	}

	size, err := w.core.StoreBackup(backup, r)
	if err != nil {
		return fmt.Errorf("StoreBackup error: %v", err)
	}

	backup.Size = size
	backup.State = types.BackupStateComplete

	job.State = types.JobStateDoneSuccess

	return w.repos.BackupRepo.Update(backup)
}
