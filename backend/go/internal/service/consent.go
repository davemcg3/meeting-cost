package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ConsentService handles cookie consent management with full auditability.
type ConsentService interface {
	// Consent management
	GetConsent(ctx context.Context, sessionID string) (*ConsentDTO, error)
	UpdateConsent(ctx context.Context, req UpdateConsentRequest) (*ConsentDTO, error)
	WithdrawConsent(ctx context.Context, sessionID string, cookieTypes []string) error

	// Cookie enforcement
	CheckCookieAllowed(ctx context.Context, sessionID string, cookieCategory string) (bool, error)
	ClassifyCookie(cookieName string) string // Returns: "necessary", "analytics", "marketing", "functional"

	// Audit and compliance
	GetConsentHistory(ctx context.Context, sessionID string, personID *uuid.UUID) ([]*ConsentDTO, error)
	ExportConsentData(ctx context.Context, personID uuid.UUID) (*ConsentExportDTO, error)

	// Policy management
	GetCurrentPolicyVersion(ctx context.Context) (string, error)

	// Syncing
	SyncConsent(ctx context.Context, sessionID string, personID uuid.UUID) error
}

type UpdateConsentRequest struct {
	SessionID         string     `json:"session_id" validate:"required"`
	PersonID          *uuid.UUID `json:"person_id"`
	AnalyticsCookies  bool       `json:"analytics_cookies"`
	MarketingCookies  bool       `json:"marketing_cookies"`
	FunctionalCookies bool       `json:"functional_cookies"`
	IPAddress         string     `json:"-"` // Set from request context
	UserAgent         string     `json:"-"` // Set from request context
}

type ConsentDTO struct {
	ID                uuid.UUID  `json:"id"`
	PersonID          *uuid.UUID `json:"person_id,omitempty"`
	SessionID         string     `json:"session_id"`
	NecessaryCookies  bool       `json:"necessary_cookies"`
	AnalyticsCookies  bool       `json:"analytics_cookies"`
	MarketingCookies  bool       `json:"marketing_cookies"`
	FunctionalCookies bool       `json:"functional_cookies"`
	ConsentVersion    string     `json:"consent_version"`
	ConsentDate       time.Time  `json:"consent_date"`
	PreviousConsentID *uuid.UUID `json:"previous_consent_id,omitempty"`
}

type ConsentExportDTO struct {
	PersonID   uuid.UUID    `json:"person_id"`
	Consents   []ConsentDTO `json:"consents"`
	ExportDate time.Time    `json:"export_date"`
}
