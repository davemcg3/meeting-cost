package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Increment represents a time slice of a meeting with specific attendee count and wage settings.
type Increment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Foreign key
	MeetingID uuid.UUID `gorm:"type:uuid;not null;index:idx_increment_meeting" json:"meeting_id"`

	// Time boundaries
	StartTime time.Time `gorm:"not null;index:idx_increment_time" json:"start_time"`
	StopTime  time.Time `gorm:"not null;index:idx_increment_time" json:"stop_time"`

	// Increment details
	AttendeeCount int     `gorm:"not null" json:"attendee_count"`
	AverageWage   float64 `gorm:"type:decimal(10,2);not null" json:"average_wage"`

	// Computed fields
	ElapsedTime int     `gorm:"not null" json:"elapsed_time"` // seconds
	Cost        float64 `gorm:"type:decimal(12,2);not null" json:"cost"`
	TotalCost   float64 `gorm:"type:decimal(12,2);not null" json:"total_cost"` // Running total at end of increment

	// Purpose (copied from meeting at increment creation)
	Purpose string `gorm:"type:text" json:"purpose"`

	// Relationships
	Meeting Meeting `gorm:"foreignKey:MeetingID" json:"-"`
}

// TableName overrides the table name.
func (Increment) TableName() string {
	return "increments"
}

// BeforeCreate ensures UUID is set if not already.
func (i *Increment) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
