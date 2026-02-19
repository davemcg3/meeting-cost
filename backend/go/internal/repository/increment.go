package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
)

// IncrementRepository handles operations for Increment entities.
type IncrementRepository interface {
	// Create
	Create(ctx context.Context, increment *models.Increment) error
	CreateBatch(ctx context.Context, increments []*models.Increment) error

	// Read
	GetByID(ctx context.Context, id uuid.UUID) (*models.Increment, error)
	GetByMeeting(ctx context.Context, meetingID uuid.UUID) ([]*models.Increment, error)

	// Update
	Update(ctx context.Context, increment *models.Increment) error

	// Delete (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByMeeting(ctx context.Context, meetingID uuid.UUID) error
}

