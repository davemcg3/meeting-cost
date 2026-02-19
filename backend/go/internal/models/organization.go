package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Organization represents an organization that can have multiple people and meetings.
type Organization struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Identity
	Name        string `gorm:"not null" json:"name"`
	Slug        string `gorm:"uniqueIndex:idx_org_slug;not null" json:"slug"` // URL-friendly identifier
	Description string `gorm:"type:text" json:"description"`

	// Default wage settings
	DefaultWage     float64 `gorm:"type:decimal(10,2);default:0" json:"default_wage"` // Default hourly wage
	UseBlendedWage bool    `gorm:"default:false" json:"use_blended_wage"`              // Use blended wage instead of individual

	// Settings - flexible storage
	Settings datatypes.JSON `gorm:"type:jsonb" json:"settings,omitempty"`
}

// TableName overrides the table name.
func (Organization) TableName() string {
	return "organizations"
}

// BeforeCreate ensures UUID is set if not already.
func (o *Organization) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
