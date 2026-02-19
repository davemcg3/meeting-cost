package service

import (
	"context"

	"github.com/google/uuid"
)

// AuditLogService handles creating audit logs.
type AuditLogService interface {
	Log(ctx context.Context, params LogParams) error
}

// LogParams contains data for creating an audit log.
type LogParams struct {
	PersonID       *uuid.UUID
	OrganizationID *uuid.UUID
	Action         string
	ResourceType   string
	ResourceID     uuid.UUID
	Details        map[string]interface{}
	IPAddress      string
	UserAgent      string
}
