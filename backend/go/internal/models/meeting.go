package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Meeting represents a meeting with its metadata and associated increments.
type Meeting struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Organization scope
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;index:idx_meeting_org" json:"organization_id"`

	// Meeting metadata
	Purpose   string     `gorm:"type:text" json:"purpose"`
	StartedAt *time.Time `json:"started_at,omitempty"` // Null if not started
	StoppedAt *time.Time `json:"stopped_at,omitempty"` // Null if still running
	IsActive  bool       `gorm:"default:false;index:idx_meeting_active" json:"is_active"`

	// Deduplication
	ExternalID        string `gorm:"index:idx_meeting_external" json:"external_id,omitempty"`         // Zoom/Teams/Slack meeting ID
	ExternalType      string `gorm:"type:varchar(50)" json:"external_type,omitempty"`                 // "zoom", "teams", "slack", "google"
	DeduplicationHash string `gorm:"index:idx_meeting_dedup" json:"deduplication_hash,omitempty"`       // Hash for deduplication

	// Creator
	CreatedByID uuid.UUID `gorm:"type:uuid;not null;index" json:"created_by_id"`

	// Computed fields (cached for performance)
	TotalCost     float64 `gorm:"type:decimal(12,2);default:0" json:"total_cost"`
	TotalDuration int     `gorm:"default:0" json:"total_duration"` // seconds
	MaxAttendees  int     `gorm:"default:0" json:"max_attendees"`

	// Relationships (for preloading)
	Organization Organization        `gorm:"foreignKey:OrganizationID" json:"-"`
	CreatedBy    Person              `gorm:"foreignKey:CreatedByID" json:"-"`
	Increments   []Increment         `gorm:"foreignKey:MeetingID" json:"-"`
	Participants []MeetingParticipant `gorm:"foreignKey:MeetingID" json:"-"`
}

// TableName overrides the table name.
func (Meeting) TableName() string {
	return "meetings"
}

// BeforeCreate ensures UUID is set if not already.
func (m *Meeting) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
