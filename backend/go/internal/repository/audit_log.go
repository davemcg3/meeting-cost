package repository

import (
	"context"

	"github.com/yourorg/meeting-cost/backend/go/internal/models"
)

// AuditLogRepository handles all database operations for AuditLog entities.
type AuditLogRepository interface {
	Create(ctx context.Context, auditLog *models.AuditLog) error
}
