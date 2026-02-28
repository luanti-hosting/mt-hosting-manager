package db

import (
	"mt-hosting-manager/types"

	"gorm.io/gorm"
)

type ImageVersionRepository struct {
	g *gorm.DB
}

func (r *ImageVersionRepository) Insert(c *types.ImageVersion) error {
	return r.g.Create(c).Error
}

func (r *ImageVersionRepository) Update(n *types.ImageVersion) error {
	return r.g.Model(n).Updates(n).Error
}

func (r *ImageVersionRepository) GetByName(name types.ImageName) (*types.ImageVersion, error) {
	return FindSingle[types.ImageVersion](r.g.Where(types.ImageVersion{Name: name}))
}

func (r *ImageVersionRepository) GetAll() ([]*types.ImageVersion, error) {
	return FindMulti[types.ImageVersion](r.g.Where(types.ImageVersion{}))
}
