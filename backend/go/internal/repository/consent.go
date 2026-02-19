package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
)

// ConsentRepository handles cookie consent database operations.
type ConsentRepository interface {
	// Create
	Create(ctx context.Context, consent *models.CookieConsent) error

	// Read
	GetByID(ctx context.Context, id uuid.UUID) (*models.CookieConsent, error)
	GetCurrentBySession(ctx context.Context, sessionID string) (*models.CookieConsent, error)
	GetCurrentByPerson(ctx context.Context, personID uuid.UUID) (*models.CookieConsent, error)
	GetHistoryBySession(ctx context.Context, sessionID string) ([]*models.CookieConsent, error)
	GetHistoryByPerson(ctx context.Context, personID uuid.UUID) ([]*models.CookieConsent, error)

	// Update
	Update(ctx context.Context, consent *models.CookieConsent) error

	// Delete (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error
}

