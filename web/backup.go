package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/types"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *Api) CreateBackup(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	b := &types.Backup{}
	err := json.NewDecoder(r.Body).Decode(b)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	// assign current user's id
	b.UserID = c.UserID

	// validate given id's
	mtserver, err := a.repos.MinetestServerRepo.GetByID(b.MinetestServerID)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if mtserver == nil {
		SendError(w, 404, fmt.Errorf("server not found: '%s'", b.MinetestServerID))
		return
	}
	node, err := a.repos.UserNodeRepo.GetByID(mtserver.UserNodeID)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if node == nil {
		SendError(w, 404, fmt.Errorf("node not found: '%s'", mtserver.UserNodeID))
		return
	}
	if node.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, 403, fmt.Errorf("not your node and not an admin"))
		return
	}

	b.State = types.BackupStateCreated
	b.Passphrase = core.RandStringRunes(64)
	b.ID = uuid.NewString()
	b.Size = 0
	b.Created = time.Now().Unix()
	b.UserID = c.UserID

	err = a.repos.BackupRepo.Insert(b)
	if err != nil {
		SendError(w, 500, fmt.Errorf("backup insert error: %v", err))
		return
	}

	// create job
	job := types.BackupServerJob(node, mtserver, b)
	err = a.repos.JobRepo.Insert(job)
	if err != nil {
		SendError(w, 500, fmt.Errorf("backup-job insert error: %v", err))
		return
	}

	Send(w, b, nil)
}

func (a *Api) GetBackups(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	if c.Role == types.UserRoleAdmin {
		// return all
		list, err := a.repos.BackupRepo.GetAll()
		Send(w, list, err)
	} else {
		// return only the user's backups
		list, err := a.repos.BackupRepo.GetByUserID(c.UserID)
		Send(w, list, err)
	}
}

func (a *Api) RemoveBackup(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	b, err := a.repos.BackupRepo.GetByID(vars["id"])
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if b == nil {
		SendError(w, 404, fmt.Errorf("backup not found"))
		return
	}
	if b.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	err = a.core.RemoveBackup(b)
	Send(w, true, err)
}

func (a *Api) GetBackup(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	b, err := a.repos.BackupRepo.GetByID(vars["id"])
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if b == nil {
		SendError(w, 404, fmt.Errorf("backup not found"))
		return
	}

	if b.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	Send(w, b, nil)
}

func (a *Api) UpdateBackup(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	b, err := a.repos.BackupRepo.GetByID(vars["id"])
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if b == nil {
		SendError(w, 404, fmt.Errorf("backup not found"))
		return
	}

	if b.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	updated_backup := &types.Backup{}
	err = json.NewDecoder(r.Body).Decode(updated_backup)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json error: %v", err))
		return
	}

	// update allowed fields
	b.Comment = updated_backup.Comment
	b.Expires = updated_backup.Expires

	err = a.repos.BackupRepo.Update(b)
	Send(w, b, err)
}

func (a *Api) DownloadBackup(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	b, err := a.repos.BackupRepo.GetByID(vars["id"])
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if b == nil {
		SendError(w, 404, fmt.Errorf("backup not found"))
		return
	}

	if b.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", b.ID))
	w.Header().Set("Content-Type", "application/gzip")

	err = a.core.StreamBackup(b, w)
	if err != nil {
		SendError(w, 500, err)
	}
}

func (a *Api) GetBackupJob(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	b, err := a.repos.BackupRepo.GetByID(vars["id"])
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if b == nil {
		SendError(w, 404, fmt.Errorf("backup not found"))
		return
	}

	if b.UserID != c.UserID && c.Role != types.UserRoleAdmin {
		SendError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	job, err := a.repos.JobRepo.GetLatestByBackupID(b.ID)
	Send(w, job, err)
}
