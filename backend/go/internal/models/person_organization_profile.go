package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PersonOrganizationProfile is the join table between Person and Organization,
// including wage information and membership status.
type PersonOrganizationProfile struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign keys
	PersonID       uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_person_org_unique" json:"person_id"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_person_org_unique" json:"organization_id"`

	// Membership
	IsActive bool      `gorm:"default:true" json:"is_active"`
	JoinedAt time.Time `gorm:"not null" json:"joined_at"`
	LeftAt   *time.Time `json:"left_at,omitempty"` // Null if still active

	// Wage information (nullable; uses org default if null)
	HourlyWage    *float64   `gorm:"type:decimal(10,2)" json:"hourly_wage,omitempty"`
	WageUpdatedAt *time.Time `json:"wage_updated_at,omitempty"`

	// External IDs for meeting integration (Zoom, Teams, Slack, etc.)
	ExternalIDs datatypes.JSON `gorm:"type:jsonb" json:"external_ids,omitempty"`

	// Relationships (for preloading)
	Person       Person       `gorm:"foreignKey:PersonID" json:"-"`
	Organization Organization `gorm:"foreignKey:OrganizationID" json:"-"`
}

// TableName overrides the table name.
func (PersonOrganizationProfile) TableName() string {
	return "person_organization_profiles"
}

// BeforeCreate ensures UUID is set if not already.
func (p *PersonOrganizationProfile) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
