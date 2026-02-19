package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthMethod represents an authentication method linked to a person (OAuth providers, email/password, etc.).
type AuthMethod struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Person association
	PersonID uuid.UUID `gorm:"type:uuid;not null;index:idx_auth_method_person" json:"person_id"`

	// Auth method details
	Provider   string `gorm:"type:varchar(50);not null;uniqueIndex:idx_auth_method_provider" json:"provider"`   // "email", "oauth_zoom", "oauth_google", etc.
	ProviderID string `gorm:"not null;uniqueIndex:idx_auth_method_provider" json:"provider_id"` // External provider's user ID
	Email      string `gorm:"index:idx_auth_method_email" json:"email"`   // Email from provider

	// OAuth tokens (stored encrypted at application level)
	AccessToken  string     `gorm:"type:text" json:"-"`
	RefreshToken string     `gorm:"type:text" json:"-"`
	TokenExpiry  *time.Time `json:"token_expiry,omitempty"`

	// Password (hashed, only for email provider)
	PasswordHash string `gorm:"type:varchar(255)" json:"-"`

	// Verification
	EmailVerified bool       `gorm:"default:false" json:"email_verified"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`

	// Relationships
	Person Person `gorm:"foreignKey:PersonID" json:"-"`
}

// TableName overrides the table name.
func (AuthMethod) TableName() string {
	return "auth_methods"
}

// BeforeCreate ensures UUID is set if not already.
func (a *AuthMethod) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
