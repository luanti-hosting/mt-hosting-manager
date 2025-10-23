package db

import (
	"mt-hosting-manager/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	g *gorm.DB
}

func (r *UserRepository) Insert(u *types.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	if u.Created == 0 {
		u.Created = time.Now().Unix()
	}
	return r.g.Create(u).Error
}

func (r *UserRepository) Update(u *types.User) error {
	return r.g.Model(u).Select("*").Updates(u).Error
}

func (r *UserRepository) GetByID(id string) (*types.User, error) {
	var list []*types.User
	err := r.g.Where(types.User{ID: id}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *UserRepository) GetByName(name string) (*types.User, error) {
	var list []*types.User
	err := r.g.Where(types.User{Name: name}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *UserRepository) GetByTypeAndExternalID(t types.UserType, external_id string) (*types.User, error) {
	var list []*types.User
	err := r.g.Where(types.User{Type: t, ExternalID: external_id}).Limit(1).Find(&list).Error
	if len(list) == 0 {
		return nil, err
	}
	return list[0], err
}

func (r *UserRepository) GetAllByRole(role types.UserRole) ([]*types.User, error) {
	return FindMulti[types.User](r.g.Where(types.User{Role: role}).Order("lastlogin DESC"))
}

func (r *UserRepository) GetAll() ([]*types.User, error) {
	return FindMulti[types.User](r.g.Where(types.User{}).Order("lastlogin DESC"))
}

func (r *UserRepository) CountAll() (int64, error) {
	var c int64
	return c, r.g.Model(types.User{}).Count(&c).Error
}

func (r *UserRepository) Delete(user_id string) error {
	return r.g.Delete(types.User{ID: user_id}).Error
}

func (r *UserRepository) DeleteAll() error {
	return r.g.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(types.User{}).Error
}

func (r *UserRepository) AddBalance(user_id string, eurocents int64) error {
	return r.g.Exec("update public.user set balance = balance + $1 where id = $2", eurocents, user_id).Error
}

func (r *UserRepository) SubtractBalance(user_id string, eurocents int64) error {
	return r.g.Exec("update public.user set balance = balance - $1 where id = $2", eurocents, user_id).Error
}

func (r *UserRepository) Search(s *types.UserSearch) ([]*types.User, error) {
	q := r.g

	if s.NameLike != nil {
		q = q.Where("name like ?", *s.NameLike)
	}

	if s.UserID != nil {
		q = q.Where(types.User{ID: *s.UserID})
	}

	if s.Limit != nil && *s.Limit > 0 && *s.Limit < 100 {
		q = q.Limit(*s.Limit)
	} else {
		q = q.Limit(100)
	}

	var list []*types.User
	err := q.Order("lastlogin DESC").Find(&list).Error
	return list, err
}
