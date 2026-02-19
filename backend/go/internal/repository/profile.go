package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
)

// PersonOrganizationProfileRepository handles operations for the Person-Organization relationship.
type PersonOrganizationProfileRepository interface {
	// Create
	Create(ctx context.Context, profile *models.PersonOrganizationProfile) error

	// Read
	GetByID(ctx context.Context, id uuid.UUID) (*models.PersonOrganizationProfile, error)
	GetByPersonAndOrg(ctx context.Context, personID, orgID uuid.UUID) (*models.PersonOrganizationProfile, error)
	GetByPerson(ctx context.Context, personID uuid.UUID) ([]*models.PersonOrganizationProfile, error)
	GetByOrganization(ctx context.Context, orgID uuid.UUID, activeOnly bool) ([]*models.PersonOrganizationProfile, error)

	// Update
	Update(ctx context.Context, profile *models.PersonOrganizationProfile) error
	UpdateWage(ctx context.Context, personID, orgID uuid.UUID, wage float64) error

	// Membership
	Activate(ctx context.Context, personID, orgID uuid.UUID) error
	Deactivate(ctx context.Context, personID, orgID uuid.UUID) error

	// Delete (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error
}

