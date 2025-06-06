package worker

import (
	"errors"
	"fmt"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (w *Worker) ExecuteJob(tx *gorm.DB, job *types.Job) {
	logrus.WithFields(job.LogrusFields()).Debug("Executing job")
	w.wg.Add(1)
	defer w.wg.Done()

	var err error
	var executor = executors[job.Type]
	if executor == nil {
		err = errors.New("type not implemented")
	} else {
		err = executor(&JobContext{
			tx:  tx,
			w:   w,
			job: job,
		})
	}

	if err != nil {
		// job failed
		job.State = types.JobStateDoneFailure
		job.Message = err.Error()
		job.Finished = time.Now().Unix()

		fields := job.LogrusFields()
		fields["err"] = err
		logrus.WithFields(fields).Error("job failed")

		job_url := fmt.Sprintf("%s/#/jobs", w.cfg.BaseURL)
		notify.Send(&notify.NtfyNotification{
			Title:    fmt.Sprintf("Job failed: %s", job.Type),
			Message:  fmt.Sprintf("Type: %s, ID %s, progress %.2f, message: '%s'", job.Type, job.ID, job.ProgressPercent, job.Message),
			Priority: 3,
			Click:    &job_url,
			Tags:     []string{"arrow_forward", "warning"},
		}, true)

	} else if job.State != types.JobStateRunning {
		// done or errored
		job.Finished = time.Now().Unix()
	}

	err = w.repos.JobRepo.UpdateWithTx(tx, job)
	if err != nil {
		fields := job.LogrusFields()
		fields["err"] = err
		logrus.WithFields(fields).Error("job update failed")
		return
	}

	logrus.WithFields(job.LogrusFields()).Debug("Job step finished")
}
