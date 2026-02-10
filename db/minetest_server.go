package db

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MinetestServerRepository struct {
	g *gorm.DB
}

func (r *MinetestServerRepository) Insert(n *types.MinetestServer) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.g.Create(n).Error
}

func (r *MinetestServerRepository) Update(n *types.MinetestServer) error {
	return r.g.Model(n).Updates(n).Error
}

func (r *MinetestServerRepository) GetByID(id string) (*types.MinetestServer, error) {
	return FindSingle[types.MinetestServer](r.g.Where(types.MinetestServer{ID: id}))
}

func (r *MinetestServerRepository) GetAll() ([]*types.MinetestServer, error) {
	return FindMulti[types.MinetestServer](r.g.Where(types.MinetestServer{}))
}

func (r *MinetestServerRepository) Delete(id string) error {
	return r.g.Delete(types.MinetestServer{ID: id}).Error
}

func (r *MinetestServerRepository) Search(search *types.MinetestServerSearch) ([]*types.MinetestServer, error) {
	q := r.g

	if search.ID != nil {
		q = q.Where(types.MinetestServer{ID: *search.ID})
	}

	if search.UserID != nil {
		q = q.Where("user_node_id in (select id from user_node where user_id = ?)", *search.UserID)
	}

	if search.NodeID != nil {
		q = q.Where(types.MinetestServer{UserNodeID: *search.NodeID})
	}

	if search.State != nil {
		q = q.Where(types.MinetestServer{State: *search.State})
	}

	return FindMulti[types.MinetestServer](q)
}
