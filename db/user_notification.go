package db

import (
	"mt-hosting-manager/types"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserNotificationRepository struct {
	g *gorm.DB
}

func (r *UserNotificationRepository) Insert(n *types.UserNotification) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	return r.g.Create(n).Error
}

func (r *UserNotificationRepository) GetUnread(UserID string) ([]*types.UserNotification, error) {
	return FindMulti[types.UserNotification](r.g.Where("seen = false and user_id = ?", UserID))
}

func (r *UserNotificationRepository) GetByID(id string) (*types.UserNotification, error) {
	return FindSingle[types.UserNotification](r.g.Where(types.UserNotification{ID: id}))
}

func (r *UserNotificationRepository) MarkRead(id string) error {
	return r.g.Model(types.UserNotification{}).Where(types.UserNotification{ID: id}).Update("seen", true).Error
}
