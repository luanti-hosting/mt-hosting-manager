package db

import (
	"gorm.io/gorm"
)

type Repositories struct {
	UserRepo               *UserRepository
	UserSettingRepo        *UserSettingRepository
	UserNotificationRepo   *UserNotificationRepository
	NodeTypeRepo           *NodeTypeRepository
	UserNodeRepo           *UserNodeRepository
	MinetestServerRepo     *MinetestServerRepository
	JobRepo                *JobRepository
	PaymentTransactionRepo *PaymentTransactionRepository
	AuditLogRepo           *AuditLogRepository
	ServiceTicketRepo      *ServiceTicketRepository
	BackupRepo             *BackupRepository
	ExchangeRateRepo       *ExchangeRateRepository
	CouponRepo             *CouponRepository
	Lock                   *DBLock
	g                      *gorm.DB
}

func NewRepositories(g *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo:               &UserRepository{g: g},
		UserSettingRepo:        &UserSettingRepository{g: g},
		UserNotificationRepo:   &UserNotificationRepository{g: g},
		NodeTypeRepo:           &NodeTypeRepository{g: g},
		UserNodeRepo:           &UserNodeRepository{g: g},
		MinetestServerRepo:     &MinetestServerRepository{g: g},
		JobRepo:                &JobRepository{g: g},
		PaymentTransactionRepo: &PaymentTransactionRepository{g: g},
		AuditLogRepo:           &AuditLogRepository{g: g},
		ServiceTicketRepo:      &ServiceTicketRepository{g: g},
		BackupRepo:             &BackupRepository{g: g},
		ExchangeRateRepo:       &ExchangeRateRepository{g: g},
		CouponRepo:             &CouponRepository{g: g},
		Lock:                   &DBLock{g: g},
		g:                      g,
	}
}

func (r *Repositories) Gorm() *gorm.DB {
	return r.g
}
