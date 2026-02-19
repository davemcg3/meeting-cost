package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AuditLog tracks important actions for security and compliance.
// No DeletedAt - audit logs are immutable and never soft-deleted.
type AuditLog struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"index:idx_audit_created_at" json:"created_at"`

	// Actor
	PersonID       *uuid.UUID `gorm:"type:uuid;index:idx_audit_person" json:"person_id,omitempty"`
	OrganizationID *uuid.UUID `gorm:"type:uuid;index:idx_audit_org" json:"organization_id,omitempty"`

	// Action details
	Action       string `gorm:"type:varchar(100);not null" json:"action"`       // "create", "update", "delete", "login", "logout"
	ResourceType string `gorm:"type:varchar(50);not null" json:"resource_type"` // "meeting", "organization", "person"
	ResourceID   uuid.UUID `gorm:"type:uuid;index:idx_audit_resource" json:"resource_id"`

	// Details
	Details   datatypes.JSON `gorm:"type:jsonb" json:"details,omitempty"`
	IPAddress string         `json:"ip_address,omitempty"`
	UserAgent string         `json:"user_agent,omitempty"`
}

// TableName overrides the table name.
func (AuditLog) TableName() string {
	return "audit_logs"
}
