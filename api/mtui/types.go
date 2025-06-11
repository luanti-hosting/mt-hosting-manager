package mtui

type BackupRestoreType string

const (
	BackupJob  BackupRestoreType = "backup"
	RestoreJob BackupRestoreType = "restore"
)

type BackupRestoreJobState string

const (
	BackupRestoreJobRunning BackupRestoreJobState = "running"
	BackupRestoreJobSuccess BackupRestoreJobState = "success"
	BackupRestoreJobFailure BackupRestoreJobState = "failure"
)

// new job
type CreateBackupRestoreJob struct {
	Type BackupRestoreType `json:"type"`

	Endpoint  string `json:"endpoint"`
	KeyID     string `json:"key_id"`
	AccessKey string `json:"access_key"`
	Bucket    string `json:"bucket"`

	FileKey  string `json:"file_key"`
	Filename string `json:"filename"`
}

// current job info
type BackupRestoreInfo struct {
	Type            BackupRestoreType     `json:"type"`
	ProgressPercent float64               `json:"progress_percent"`
	Message         string                `json:"message"`
	State           BackupRestoreJobState `json:"state"`
}

type Stats struct {
	Uptime      int     `json:"uptime"`
	MaxLag      float64 `json:"max_lag"`
	TimeOfDay   float64 `json:"time_of_day"`
	PlayerCount int     `json:"player_count"`
	Maintenance bool    `json:"maintenance"`
}
