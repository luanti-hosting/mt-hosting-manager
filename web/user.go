package web

import (
	"encoding/json"
	"mt-hosting-manager/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) SearchUser(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	s := &types.UserSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	list, err := a.repos.UserRepo.Search(s)
	for _, u := range list {
		u.RemoveSensitiveFields()
	}
	Send(w, list, err)
}

func (a *Api) GetUsers(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	list, err := a.repos.UserRepo.GetAll()
	Send(w, list, err)
}

func (a *Api) GetUserByID(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	vars := mux.Vars(r)
	user, err := a.repos.UserRepo.GetByID(vars["id"])

	if c.Role != types.UserRoleAdmin {
		// non-admin, return stripped-down object
		user = &types.User{
			ID:   user.ID,
			Name: user.Name,
			Role: user.Role,
		}
	}

	Send(w, user, err)
}

func (a *Api) SaveUser(w http.ResponseWriter, r *http.Request, c *types.Claims) {
	updated_user := &types.User{}
	err := json.NewDecoder(r.Body).Decode(updated_user)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	user, err := a.repos.UserRepo.GetByID(updated_user.ID)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	user.Balance = updated_user.Balance
	user.Role = updated_user.Role
	user.State = updated_user.State
	user.Hash = updated_user.Hash
	user.Type = updated_user.Type
	user.Name = updated_user.Name

	err = a.repos.UserRepo.Update(user)
	Send(w, user, err)
}
