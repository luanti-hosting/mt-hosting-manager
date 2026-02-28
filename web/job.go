package web

import (
	"fmt"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (a *Api) GetJobs(w http.ResponseWriter, r *http.Request) {
	list, err := a.repos.JobRepo.GetAll()
	Send(w, list, err)
}

func (a *Api) RetryJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	job, err := a.repos.JobRepo.GetByID(id)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if job == nil {
		SendError(w, 404, fmt.Errorf("job %s not found", id))
		return
	}

	// reset state, message and next run date
	job.State = types.JobStateRunning
	job.Message = ""
	job.NextRun = time.Now().Unix()
	err = a.repos.JobRepo.Update(job)
	Send(w, job, err)
}

func (a *Api) DeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := a.repos.JobRepo.Delete(vars["id"])
	Send(w, true, err)
}
