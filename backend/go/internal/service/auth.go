package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
)

// AuthService handles authentication and authorization logic.
type AuthService interface {
	// Registration
	Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	VerifyEmail(ctx context.Context, token string) error

	// Authentication
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, token string, ipAddress, userAgent string) error
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)

	// OAuth
	OAuthLogin(ctx context.Context, provider string, code string) (*LoginResponse, error)
	OAuthCallback(ctx context.Context, provider string, state, code string) (*LoginResponse, error)
	LinkOAuthProvider(ctx context.Context, personID uuid.UUID, provider string, code string) error

	// Password management
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	ChangePassword(ctx context.Context, personID uuid.UUID, oldPassword, newPassword string) error

	// Session management
	ValidateSession(ctx context.Context, token string) (*SessionInfo, error)
	GetSessions(ctx context.Context, personID uuid.UUID) ([]*models.Session, error)
	RevokeSession(ctx context.Context, personID, sessionID uuid.UUID) error
	RevokeAllSessions(ctx context.Context, personID uuid.UUID) error
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"`
	IPAddress string `json:"-"`
	UserAgent string `json:"-"`
}

type RegisterResponse struct {
	User        *models.Person `json:"user"`
	AccessToken string         `json:"access_token"`
	ExpiresIn   int            `json:"expires_in"`
}

type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	IPAddress string `json:"-"`
	UserAgent string `json:"-"`
}

type LoginResponse struct {
	User         *models.Person `json:"user"`
	AccessToken  string         `json:"access_token"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    int            `json:"expires_in"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type SessionInfo struct {
	PersonID     uuid.UUID
	Email        string
	ExpiresAt    time.Time
	LastActivity time.Time
}
