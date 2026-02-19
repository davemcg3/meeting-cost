package gorm

import (
	"context"
	"fmt"

	"github.com/yourorg/meeting-cost/backend/go/internal/models"
	"github.com/yourorg/meeting-cost/backend/go/internal/repository"
	"gorm.io/gorm"
)

type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository creates a new GORM-based AuditLogRepository.
func NewAuditLogRepository(db *gorm.DB) repository.AuditLogRepository {
	return &auditLogRepository{
		db: db,
	}
}

func (r *auditLogRepository) Create(ctx context.Context, auditLog *models.AuditLog) error {
	if err := r.db.WithContext(ctx).Create(auditLog).Error; err != nil {
		return fmt.Errorf("creating audit log: %w", err)
	}
	return nil
}
