package web

import (
	"encoding/json"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) GetNodeTypes(w http.ResponseWriter, r *http.Request) {
	list, err := a.repos.NodeTypeRepo.GetAll()
	if r.URL.Query().Get("export") == "true" {
		w.Header().Add("Content-Disposition", "attachment; filename=\"nodetypes.json\"")
	}
	Send(w, list, err)
}

func (a *Api) GetNodeType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nt, err := a.repos.NodeTypeRepo.GetByID(vars["id"])
	Send(w, nt, err)
}

func (a *Api) CreateNodeType(w http.ResponseWriter, r *http.Request) {
	nt := &types.NodeType{}
	err := json.NewDecoder(r.Body).Decode(nt)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	err = a.repos.NodeTypeRepo.Insert(nt)
	Send(w, nt, err)
}

func (a *Api) UpdateNodeType(w http.ResponseWriter, r *http.Request) {
	nt := &types.NodeType{}
	err := json.NewDecoder(r.Body).Decode(nt)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	err = a.repos.NodeTypeRepo.Update(nt)
	Send(w, nt, err)
}

func (a *Api) DeleteNodeType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := a.repos.NodeTypeRepo.Delete(vars["id"])
	Send(w, true, err)
}
