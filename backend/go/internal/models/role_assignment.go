package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleAssignment links roles to persons within an organization.
type RoleAssignment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	RoleID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_role_assignment" json:"role_id"`
	PersonID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_role_assignment" json:"person_id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_role_assignment" json:"organization_id"`

	// Relationships (for preloading)
	Role         Role         `gorm:"foreignKey:RoleID" json:"-"`
	Person       Person       `gorm:"foreignKey:PersonID" json:"-"`
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"-"`
}

// TableName overrides the table name.
func (RoleAssignment) TableName() string {
	return "role_assignments"
}

// BeforeCreate ensures UUID is set if not already.
func (r *RoleAssignment) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
