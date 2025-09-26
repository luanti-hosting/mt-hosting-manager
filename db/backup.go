package db

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BackupRepository struct {
	g *gorm.DB
}

func (r *BackupRepository) Insert(n *types.Backup) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.g.Create(n).Error
}

func (r *BackupRepository) Update(n *types.Backup) error {
	return r.g.Model(n).Updates(n).Error
}

func (r *BackupRepository) GetByID(id string) (*types.Backup, error) {
	return FindSingle[types.Backup](r.g.Where(types.Backup{ID: id}))
}

func (r *BackupRepository) GetByState(state types.BackupState) ([]*types.Backup, error) {
	return FindMulti[types.Backup](r.g.Where(types.Backup{State: state}))
}

func (r *BackupRepository) GetByUserID(user_id string) ([]*types.Backup, error) {
	return FindMulti[types.Backup](r.g.Where(types.Backup{UserID: user_id}))
}

func (r *BackupRepository) GetAll() ([]*types.Backup, error) {
	return FindMulti[types.Backup](r.g.Where(types.Backup{}))
}

func (r *BackupRepository) Delete(id string) error {
	return r.g.Delete(types.Backup{ID: id}).Error
}
