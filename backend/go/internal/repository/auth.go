package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
)

// AuthRepository handles authentication-related database operations.
type AuthRepository interface {
	// AuthMethod operations
	CreateAuthMethod(ctx context.Context, method *models.AuthMethod) error
	GetAuthMethodByID(ctx context.Context, id uuid.UUID) (*models.AuthMethod, error)
	GetAuthMethodByProvider(ctx context.Context, provider, providerID string) (*models.AuthMethod, error)
	GetAuthMethodsByPerson(ctx context.Context, personID uuid.UUID) ([]*models.AuthMethod, error)
	UpdateAuthMethod(ctx context.Context, method *models.AuthMethod) error
	DeleteAuthMethod(ctx context.Context, id uuid.UUID) error

	// Session operations
	CreateSession(ctx context.Context, session *models.Session) error
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (*models.Session, error)
	GetSessionsByPerson(ctx context.Context, personID uuid.UUID) ([]*models.Session, error)
	UpdateSession(ctx context.Context, session *models.Session) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteExpiredSessions(ctx context.Context) error
	DeleteSessionsByPerson(ctx context.Context, personID uuid.UUID) error
}

