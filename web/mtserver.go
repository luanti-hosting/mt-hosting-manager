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

func (a *Api) GetMTServers(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	if c.Role == types.UserRoleAdmin && r.URL.Query().Get("all") == "true" {
		list, err := a.repos.MinetestServerRepo.GetAll()
		Send(w, list, err)
	} else {
		list, err := a.repos.MinetestServerRepo.Search(&types.MinetestServerSearch{UserID: &c.UserID})
		Send(w, list, err)
	}
}

func (a *Api) SearchMTServers(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	search := &types.MinetestServerSearch{}
	err := json.NewDecoder(r.Body).Decode(search)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	if c.Role != types.UserRoleAdmin {
		// fix userid
		search.UserID = &c.UserID
	}

	list, err := a.repos.MinetestServerRepo.Search(search)
	Send(w, list, err)
}

func (a *Api) GetMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, server, _, err := a.CheckedGetMTServer(id, c)
	Send(w, server, err)
}

func (a *Api) ValidateCreateMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	create_mtserver := &types.MinetestServer{}
	err := json.NewDecoder(r.Body).Decode(create_mtserver)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	node, status, err := a.CheckedGetUserNode(create_mtserver.UserNodeID, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	v, err := a.core.ValidateCreateServer(create_mtserver, node)
	Send(w, v, err)
}

func (a *Api) CreateMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	create_mtserver := &types.MinetestServer{}
	err := json.NewDecoder(r.Body).Decode(create_mtserver)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	node, status, err := a.CheckedGetUserNode(create_mtserver.UserNodeID, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	v, err := a.core.ValidateCreateServer(create_mtserver, node)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if !v.Valid {
		SendError(w, 500, fmt.Errorf("validation failed"))
		return
	}

	var backup *types.Backup
	restore_from := r.URL.Query().Get("restore_from")
	if restore_from != "" {
		// restore-point specified, resolve backup
		backup, err = a.repos.BackupRepo.GetByID(restore_from)
		if err != nil {
			SendError(w, 500, fmt.Errorf("could not fetch backup '%s': %v", restore_from, err))
			return
		}

		if c.Role != types.UserRoleAdmin && backup.UserID != c.UserID {
			SendError(w, 403, fmt.Errorf("not authorized for backup id: '%s'", backup.ID))
			return
		}
	}

	mtserver := &types.MinetestServer{
		ID:         uuid.NewString(),
		UserNodeID: node.ID,
		Name:       create_mtserver.Name,
		Admin:      create_mtserver.Admin,
		DNSName:    create_mtserver.DNSName,
		Port:       create_mtserver.Port,
		Created:    time.Now().Unix(),
		UIVersion:  "master",
		JWTKey:     core.RandSeq(16),
		State:      types.MinetestServerStateCreated,
	}
	err = a.repos.MinetestServerRepo.Insert(mtserver)
	if err != nil {
		SendError(w, 500, fmt.Errorf("server insert error: %v", err))
		return
	}

	job := types.SetupServerJob(node, mtserver, backup)
	err = a.repos.JobRepo.Insert(job)
	if err != nil {
		SendError(w, 500, fmt.Errorf("job insert error: %v", err))
		return
	}

	a.core.AddAuditLog(&types.AuditLog{
		Type:             types.AuditLogServerCreated,
		UserID:           c.UserID,
		UserNodeID:       &node.ID,
		MinetestServerID: &mtserver.ID,
	})

	Send(w, mtserver, nil)
}

func (a *Api) DeleteMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	node, mtserver, status, err := a.CheckedGetMTServer(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	job := types.RemoveServerJob(node, mtserver)
	err = a.repos.JobRepo.Insert(job)

	Send(w, true, err)
}

func (a *Api) UpdateMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, mtserver, status, err := a.CheckedGetMTServer(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	updated_mtserver := &types.MinetestServer{}
	err = json.NewDecoder(r.Body).Decode(updated_mtserver)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	// update allowed fields
	mtserver.Name = updated_mtserver.Name
	mtserver.Admin = updated_mtserver.Admin
	mtserver.UIVersion = updated_mtserver.UIVersion
	mtserver.CustomDNS = updated_mtserver.CustomDNS

	err = a.repos.MinetestServerRepo.Update(mtserver)
	Send(w, mtserver, err)
}

func (a *Api) SetupMTServer(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, mtserver, status, err := a.CheckedGetMTServer(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	// check for current execution
	latest_job, err := a.repos.JobRepo.GetLatestByMinetestServerID(mtserver.ID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("could not fetch latest job: %v", err))
		return
	}
	if latest_job != nil && latest_job.State == types.JobStateRunning {
		// already running or created
		SendError(w, 500, fmt.Errorf("job already scheduled: %s", latest_job.ID))
		return
	}

	job := types.SetupServerJob(node, mtserver, nil)
	err = a.repos.JobRepo.Insert(job)

	Send(w, job, err)
}

func (a *Api) GetLatestMTServerJob(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, mtserver, status, err := a.CheckedGetMTServer(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	job, err := a.repos.JobRepo.GetLatestByMinetestServerID(mtserver.ID)
	Send(w, job, err)
}

func (a *Api) GetMTServerStats(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, server, _, err := a.CheckedGetMTServer(id, c)
	if err != nil {
		SendError(w, 500, fmt.Errorf("get server error: %v", err))
		return
	}
	if server == nil {
		SendError(w, 404, fmt.Errorf("server not found: %s", id))
		return
	}
	if server.State != types.MinetestServerStateRunning {
		SendError(w, 500, fmt.Errorf("server is not running"))
	}

	client, err := a.core.GetMTUIClient(server)
	if err != nil {
		SendError(w, 500, fmt.Errorf("get client error: %v", err))
		return
	}

	stats, err := client.GetStats()
	if err != nil {
		SendError(w, 404, fmt.Errorf("server not ready/unresponsive: %v", err))
		return
	}

	Send(w, stats, nil)
}
