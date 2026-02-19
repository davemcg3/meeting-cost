package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MeetingParticipant tracks which people participated in a meeting (for analytics and cost calculation).
type MeetingParticipant struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	MeetingID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_meeting_participant" json:"meeting_id"`
	PersonID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_meeting_participant" json:"person_id"`

	// Participation details
	JoinedAt *time.Time `json:"joined_at,omitempty"`
	LeftAt   *time.Time `json:"left_at,omitempty"`
	Duration int        `gorm:"default:0" json:"duration"` // seconds

	// Relationships
	Meeting Meeting `gorm:"foreignKey:MeetingID" json:"-"`
	Person  Person  `gorm:"foreignKey:PersonID" json:"-"`
}

// TableName overrides the table name.
func (MeetingParticipant) TableName() string {
	return "meeting_participants"
}

// BeforeCreate ensures UUID is set if not already.
func (m *MeetingParticipant) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
