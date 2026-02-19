package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
)

// PersonRepository handles all database operations for Person entities.
type PersonRepository interface {
	// Create
	Create(ctx context.Context, person *models.Person) error

	// Read
	GetByID(ctx context.Context, id uuid.UUID) (*models.Person, error)
	GetByEmail(ctx context.Context, email string) (*models.Person, error)
	List(ctx context.Context, filters PersonFilters, pagination Pagination) ([]*models.Person, int64, error)

	// Update
	Update(ctx context.Context, person *models.Person) error

	// Delete (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// Anonymization (GDPR)
	Anonymize(ctx context.Context, id uuid.UUID) error

	// Relationships
	GetOrganizations(ctx context.Context, personID uuid.UUID) ([]*models.Organization, error)
	GetActiveOrganizations(ctx context.Context, personID uuid.UUID) ([]*models.Organization, error)
}

type PersonFilters struct {
	Email          *string
	Anonymized     *bool
	OrganizationID *uuid.UUID // Filter by organization membership
}

