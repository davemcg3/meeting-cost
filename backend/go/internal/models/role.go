package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents a role within an organization that groups permissions together.
type Role struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Identity
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Organization scope
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_role_org_name" json:"organization_id"`

	// Relationships (for preloading)
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"-"`
}

// TableName overrides the table name.
func (Role) TableName() string {
	return "roles"
}

// BeforeCreate ensures UUID is set if not already.
func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
