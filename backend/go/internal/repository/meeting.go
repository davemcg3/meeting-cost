package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
)

// MeetingRepository handles all database operations for Meeting entities.
type MeetingRepository interface {
	// Create
	Create(ctx context.Context, meeting *models.Meeting) error

	// Read
	GetByID(ctx context.Context, id uuid.UUID) (*models.Meeting, error)
	GetByExternalID(ctx context.Context, externalType, externalID string) (*models.Meeting, error)
	GetByDeduplicationHash(ctx context.Context, hash string) (*models.Meeting, error)
	List(ctx context.Context, filters MeetingFilters, pagination Pagination) ([]*models.Meeting, int64, error)

	// Update
	Update(ctx context.Context, meeting *models.Meeting) error
	Start(ctx context.Context, id uuid.UUID) error
	Stop(ctx context.Context, id uuid.UUID) error

	// Delete (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// Increments
	GetIncrements(ctx context.Context, meetingID uuid.UUID) ([]*models.Increment, error)
	AddIncrement(ctx context.Context, increment *models.Increment) error

	// Participants
	GetParticipants(ctx context.Context, meetingID uuid.UUID) ([]*models.MeetingParticipant, error)
	AddParticipant(ctx context.Context, participant *models.MeetingParticipant) error
	RemoveParticipant(ctx context.Context, meetingID, personID uuid.UUID) error
}

type MeetingFilters struct {
	OrganizationID *uuid.UUID
	CreatedByID    *uuid.UUID
	IsActive       *bool
	StartedAfter   *time.Time
	StartedBefore  *time.Time
	ExternalType   *string
	ExternalID     *string
}

