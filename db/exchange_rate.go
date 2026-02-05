package db

import (
	"mt-hosting-manager/types"

	"gorm.io/gorm"
)

type ExchangeRateRepository struct {
	g *gorm.DB
}

func (r *ExchangeRateRepository) Insert(n *types.ExchangeRate) error {
	return r.g.Create(n).Error
}

func (r *ExchangeRateRepository) Update(n *types.ExchangeRate) error {
	return r.g.Model(n).Updates(n).Error
}

func (r *ExchangeRateRepository) GetAll() ([]*types.ExchangeRate, error) {
	return FindMulti[types.ExchangeRate](r.g.Where(types.ExchangeRate{}))
}

func (r *ExchangeRateRepository) GetByCurrency(currency string) (*types.ExchangeRate, error) {
	return FindSingle[types.ExchangeRate](r.g.Where(types.ExchangeRate{Currency: currency}))
}

func (r *ExchangeRateRepository) DeleteByCurrency(currency string) error {
	return r.g.Delete(types.ExchangeRate{Currency: currency}).Error
}

func (r *ExchangeRateRepository) DeleteAll() error {
	return r.g.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(types.ExchangeRate{}).Error
}
