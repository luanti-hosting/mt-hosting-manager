package db

import (
	"mt-hosting-manager/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	g *gorm.DB
}

func (r *AuditLogRepository) Insert(l *types.AuditLog) error {
	if l.ID == "" {
		l.ID = uuid.NewString()
	}
	if l.Timestamp == 0 {
		l.Timestamp = time.Now().Unix()
	}
	return r.g.Create(l).Error
}

func (r *AuditLogRepository) Search(s *types.AuditLogSearch) ([]*types.AuditLog, error) {
	q := r.g.Where("timestamp > ?", s.FromTimestamp)
	q = q.Where("timestamp < ?", s.ToTimestamp)

	if s.Type != nil {
		q = q.Where(types.AuditLog{Type: *s.Type})
	}

	if s.UserID != nil {
		q = q.Where(types.AuditLog{UserID: *s.UserID})
	}

	if s.MinetestServerID != nil {
		q = q.Where(types.AuditLog{MinetestServerID: s.MinetestServerID})
	}

	if s.UserNodeID != nil {
		q = q.Where(types.AuditLog{UserNodeID: s.UserNodeID})
	}

	if s.BackupID != nil {
		q = q.Where(types.AuditLog{BackupID: s.BackupID})
	}

	if s.PaymentTransactionID != nil {
		q = q.Where(types.AuditLog{PaymentTransactionID: s.PaymentTransactionID})
	}

	q = q.Order("timestamp desc").Limit(1000)

	return FindMulti[types.AuditLog](q)
}
