package web

import (
	"encoding/json"
	"fmt"
	"mt-hosting-manager/core"
	"mt-hosting-manager/notify"
	"mt-hosting-manager/types"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *Api) GetNodes(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	if c.Role == types.UserRoleAdmin && r.URL.Query().Get("all") == "true" {
		list, err := a.repos.UserNodeRepo.GetAll()
		Send(w, list, err)
	} else {
		list, err := a.repos.UserNodeRepo.Search(&types.UserNodeSearch{UserID: &c.UserID})
		Send(w, list, err)
	}
}

func (a *Api) SearchNodes(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	search := &types.UserNodeSearch{}
	err := json.NewDecoder(r.Body).Decode(search)
	if err != nil {
		SendError(w, 500, fmt.Errorf("json error: %v", err))
		return
	}

	if c.Role != types.UserRoleAdmin {
		// fix userid
		search.UserID = &c.UserID
	}

	list, err := a.repos.UserNodeRepo.Search(search)
	Send(w, list, err)
}

func (a *Api) GetNode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, _, err := a.CheckedGetUserNode(id, c)
	Send(w, node, err)
}

func (a *Api) GetNodeServers(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, _, err := a.CheckedGetUserNode(id, c)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	servers, err := a.repos.MinetestServerRepo.Search(&types.MinetestServerSearch{NodeID: &node.ID})
	Send(w, servers, err)
}

func (a *Api) GetLatestNodeJob(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, _, err := a.CheckedGetUserNode(id, c)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	job, err := a.repos.JobRepo.GetLatestByUserNodeID(node.ID)
	Send(w, job, err)
}

func (a *Api) GetNodeStats(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, status, err := a.CheckedGetUserNode(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	metrics := &core.NodeExporterMetrics{}

	if node.State != types.UserNodeStateRunning {
		SendError(w, 500, fmt.Errorf("node not in running state"))
		return
	}

	metrics, err = core.FetchMetrics(node.IPv4)
	Send(w, metrics, err)
}

func (a *Api) UpdateNode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, status, err := a.CheckedGetUserNode(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	updated_node := &types.UserNode{}
	err = json.NewDecoder(r.Body).Decode(updated_node)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	// update allowed fields
	node.Alias = updated_node.Alias

	err = a.repos.UserNodeRepo.Update(node)
	Send(w, node, err)
}

func (a *Api) DeleteNode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	node, status, err := a.CheckedGetUserNode(id, c)
	if err != nil {
		SendError(w, status, err)
		return
	}

	job := types.RemoveNodeJob(node)
	err = a.repos.JobRepo.Insert(job)
	if err != nil {
		SendError(w, 500, fmt.Errorf("job insert error: %v", err))
		return
	}

	Send(w, true, nil)
}

func (a *Api) CreateNode(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	create_node := &types.UserNode{}
	err := json.NewDecoder(r.Body).Decode(create_node)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	// prevent provisioning of multiple nodes at once
	provision_state := types.UserNodeStateProvisioning
	current_usernodes, err := a.repos.UserNodeRepo.Search(&types.UserNodeSearch{
		UserID: &c.UserID,
		State:  &provision_state,
	})
	if err != nil {
		SendError(w, 500, fmt.Errorf("current usernodes fetch error: %v", err))
		return
	}
	if len(current_usernodes) > 0 {
		SendError(w, 405, fmt.Errorf("provisioning limit exceeded, only one node can be created at the same time"))
	}

	nt, err := a.repos.NodeTypeRepo.GetByID(create_node.NodeTypeID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("nodetype fetch error: %v", err))
		return
	}
	if nt == nil {
		SendError(w, 404, fmt.Errorf("nodetype not found: %s", create_node.NodeTypeID))
		return
	}
	if nt.State != types.NodeTypeStateActive {
		SendError(w, 405, fmt.Errorf("nodetype in wrong state: %s", nt.State))
		return
	}

	locations := strings.Split(nt.Locations, ",")
	if !slices.Contains(locations, create_node.Location) {
		SendError(w, 404, fmt.Errorf("location not available: %s", create_node.Location))
	}

	user, err := a.repos.UserRepo.GetByID(c.UserID)
	if err != nil {
		SendError(w, 500, fmt.Errorf("user fetch error: %v", err))
		return
	}
	if user == nil {
		SendError(w, 404, fmt.Errorf("user not found: %s", c.UserID))
		return
	}

	if user.Balance < (10 * nt.DailyCost) {
		SendError(w, 405, fmt.Errorf("remaining balance is less than the daily cost of the node-type"))
		return
	}

	randstr := core.RandStringRunes(7)

	node := &types.UserNode{
		ID:         uuid.NewString(),
		UserID:     c.UserID,
		NodeTypeID: create_node.NodeTypeID,
		Location:   create_node.Location,
		Created:    time.Now().Unix(),
		ValidUntil: time.Now().Unix(),
		State:      types.UserNodeStateCreated,
		Name:       fmt.Sprintf("node-%s-%s", a.cfg.Stage, randstr),
		Alias:      create_node.Alias,
	}
	err = a.repos.UserNodeRepo.Insert(node)
	if err != nil {
		SendError(w, 500, fmt.Errorf("node insert error: %v", err))
		return
	}

	job := types.SetupNodeJob(node)
	err = a.repos.JobRepo.Insert(job)
	if err != nil {
		SendError(w, 500, fmt.Errorf("job insert error: %v", err))
		return
	}

	a.core.AddAuditLog(&types.AuditLog{
		Type:       types.AuditLogNodeCreated,
		UserID:     c.UserID,
		UserNodeID: &node.ID,
	})

	notify.Send(&notify.NtfyNotification{
		Title:    fmt.Sprintf("Node created by %s (Type: %s)", user.Name, nt.Name),
		Message:  fmt.Sprintf("User: %s, Node-type %s, Name: %s", user.Name, nt.Name, node.Name),
		Priority: 3,
		Tags:     []string{"computer", "new"},
	}, true)

	Send(w, node, nil)
}
