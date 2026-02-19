package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CookieConsent tracks user cookie consent preferences for GDPR/CCPA compliance with full auditability.
type CookieConsent struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Person association (nullable for anonymous users)
	PersonID *uuid.UUID `gorm:"type:uuid;index:idx_cookie_consent_person" json:"person_id,omitempty"`

	// Consent tracking
	SessionID string `gorm:"type:varchar(255);not null;index:idx_cookie_consent_session" json:"session_id"` // Browser session identifier

	// Consent preferences
	NecessaryCookies  bool `gorm:"default:true" json:"necessary_cookies"`  // Always true, cannot be disabled
	AnalyticsCookies  bool `gorm:"default:false" json:"analytics_cookies"`
	MarketingCookies  bool `gorm:"default:false" json:"marketing_cookies"`
	FunctionalCookies bool `gorm:"default:false" json:"functional_cookies"`

	// Consent metadata
	ConsentVersion string    `gorm:"type:varchar(50);not null" json:"consent_version"` // Version of consent policy
	ConsentDate    time.Time `gorm:"not null;index:idx_cookie_consent_date" json:"consent_date"`
	IPAddress      string   `json:"ip_address,omitempty"`
	UserAgent      string   `json:"user_agent,omitempty"`

	// Audit trail
	PreviousConsentID *uuid.UUID `gorm:"type:uuid;index:idx_cookie_consent_previous" json:"previous_consent_id,omitempty"` // Link to previous consent
	ConsentSource     string    `gorm:"type:varchar(50)" json:"consent_source,omitempty"`                                 // "initial", "update", "withdrawal"

	// Relationships
	Person Person `gorm:"foreignKey:PersonID" json:"-"`
}

// TableName overrides the table name.
func (CookieConsent) TableName() string {
	return "cookie_consents"
}

// BeforeCreate ensures UUID is set if not already.
func (c *CookieConsent) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.Must(uuid.NewRandom())
	}
	return nil
}
