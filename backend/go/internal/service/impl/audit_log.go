package impl

import (
	"context"
	"encoding/json"

	"github.com/yourorg/meeting-cost/backend/go/internal/models"
	"github.com/yourorg/meeting-cost/backend/go/internal/repository"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
	"gorm.io/datatypes"
)

type auditLogService struct {
	auditLogRepo repository.AuditLogRepository
}

// NewAuditLogService creates a new AuditLogService implementation.
func NewAuditLogService(auditLogRepo repository.AuditLogRepository) service.AuditLogService {
	return &auditLogService{
		auditLogRepo: auditLogRepo,
	}
}

func (s *auditLogService) Log(ctx context.Context, params service.LogParams) error {
	var details datatypes.JSON
	if params.Details != nil {
		b, err := json.Marshal(params.Details)
		if err == nil {
			details = datatypes.JSON(b)
		}
	}

	auditLog := &models.AuditLog{
		PersonID:       params.PersonID,
		OrganizationID: params.OrganizationID,
		Action:         params.Action,
		ResourceType:   params.ResourceType,
		ResourceID:     params.ResourceID,
		Details:        details,
		IPAddress:      params.IPAddress,
		UserAgent:      params.UserAgent,
	}

	return s.auditLogRepo.Create(ctx, auditLog)
}
