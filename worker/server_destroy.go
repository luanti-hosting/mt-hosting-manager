package worker

import (
	"fmt"
	"mt-hosting-manager/types"
)

func (w *Worker) ServerDestroy(ctx *JobContext) error {
	job := ctx.job
	node, server, err := w.GetJobContext(job)
	if err != nil {
		return err
	}

	err = w.removeServer(node, server, true)
	if err != nil {
		return fmt.Errorf("server remove error: %v", err)
	}

	job.State = types.JobStateDoneSuccess
	return nil
}
