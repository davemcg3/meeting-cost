package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Session represents an active user session (JWT token tracking).
type Session struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Person association
	PersonID uuid.UUID `gorm:"type:uuid;not null;index:idx_session_person" json:"person_id"`

	// Session details
	TokenHash    string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_session_token" json:"-"` // SHA256 of JWT
	ExpiresAt    time.Time `gorm:"not null;index:idx_session_expires" json:"expires_at"`
	LastActivity time.Time `gorm:"not null" json:"last_activity"`

	// Metadata
	UserAgent string `json:"user_agent,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`

	// Relationships
	Person Person `gorm:"foreignKey:PersonID" json:"-"`
}

// TableName overrides the table name.
func (Session) TableName() string {
	return "sessions"
}

// BeforeCreate ensures UUID is set if not already.
func (s *Session) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
