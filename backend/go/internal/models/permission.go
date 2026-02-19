package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission is a polymorphic permission attached to either a Role or a Person.
// ResourceType + ResourceID identify the owner (role or person); ResourceName + TargetResourceID + Activity define the permission.
type Permission struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Polymorphic association: "role" or "person"
	ResourceType string    `gorm:"type:varchar(50);not null;index:idx_permission_resource" json:"resource_type"`
	ResourceID   uuid.UUID `gorm:"type:uuid;not null;index:idx_permission_resource" json:"resource_id"`

	// Permission details
	ResourceName    string     `gorm:"not null" json:"resource_name"`                       // e.g., "meeting", "organization", "person"
	TargetResourceID *uuid.UUID `gorm:"type:uuid;index:idx_permission_org" json:"target_resource_id,omitempty"` // Specific resource ID, null for all
	Activity        string     `gorm:"type:varchar(20);not null" json:"activity"`           // "create", "read", "update", "delete"
	Allowed         bool       `gorm:"default:true" json:"allowed"`

	// Organization scope
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index:idx_permission_org" json:"organization_id"`

	// Relationships (for preloading; use explicit FK based on ResourceType in application code)
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"-"`
}

// TableName overrides the table name.
func (Permission) TableName() string {
	return "permissions"
}

// BeforeCreate ensures UUID is set if not already.
func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
