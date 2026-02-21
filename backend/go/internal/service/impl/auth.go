package impl

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/auth"
	"github.com/yourorg/meeting-cost/backend/go/internal/logger"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
	"github.com/yourorg/meeting-cost/backend/go/internal/repository"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
)

/*
Because you've used a struct for tokenmanager in jwt.go
tokenmanager not testable
DI inside of go
accept interface return struct
refactor
have ai look at the rest of them for consistency

[claude.md](http://claude.md/) accept interfaces return structs
*/
type authService struct {
	personRepo      repository.PersonRepository
	authRepo        repository.AuthRepository
	tokenManager    *auth.TokenManager
	auditLogService service.AuditLogService
	logger          logger.Logger
}

// NewAuthService creates a new AuthService implementation.
func NewAuthService(
	personRepo repository.PersonRepository,
	authRepo repository.AuthRepository,
	tokenManager *auth.TokenManager,
	auditLogService service.AuditLogService,
	logger logger.Logger,
) service.AuthService {
	return &authService{
		personRepo:      personRepo,
		authRepo:        authRepo,
		tokenManager:    tokenManager,
		auditLogService: auditLogService,
		logger:          logger,
	}
}

func (s *authService) Register(ctx context.Context, req service.RegisterRequest) (*service.RegisterResponse, error) {
	// 1. Check if person exists
	existing, _ := s.personRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// 2. Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	// 3. Create Person
	person := &models.Person{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
	if err := s.personRepo.Create(ctx, person); err != nil {
		return nil, fmt.Errorf("creating person: %w", err)
	}

	// 4. Create AuthMethod (Email/Password)
	authMethod := &models.AuthMethod{
		PersonID:     person.ID,
		Provider:     "email",
		ProviderID:   req.Email,
		PasswordHash: hashedPassword,
		Email:        req.Email,
	}
	if err := s.authRepo.CreateAuthMethod(ctx, authMethod); err != nil {
		return nil, fmt.Errorf("creating auth method: %w", err)
	}

	// 5. Generate Initial Token Pair
	tokens, err := s.tokenManager.GenerateTokenPair(person.ID, person.Email)
	if err != nil {
		return nil, fmt.Errorf("generating tokens: %w", err)
	}

	// 6. Create Session
	session := &models.Session{
		PersonID:  person.ID,
		TokenHash: s.hashToken(tokens.AccessToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // Match refresh token expiry
	}
	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		s.logger.Error("failed to create session after registration", "error", err)
	}

	// Audit Log
	_ = s.auditLogService.Log(ctx, service.LogParams{
		PersonID:     &person.ID,
		Action:       "register",
		ResourceType: "person",
		ResourceID:   person.ID,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
	})

	return &service.RegisterResponse{
		User:        person,
		AccessToken: tokens.AccessToken,
		ExpiresIn:   int(tokens.ExpiresIn),
	}, nil
}

func (s *authService) Login(ctx context.Context, req service.LoginRequest) (*service.LoginResponse, error) {
	// 1. Get AuthMethod by email
	// Note: We might need a repo method to get auth method by provider and email
	// or search by person email.
	person, err := s.personRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	methods, err := s.authRepo.GetAuthMethodsByPerson(ctx, person.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	var emailMethod *models.AuthMethod
	for _, m := range methods {
		if m.Provider == "email" {
			emailMethod = m
			break
		}
	}

	if emailMethod == nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 2. Verify password
	if !auth.CheckPasswordHash(req.Password, emailMethod.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 3. Generate tokens
	tokens, err := s.tokenManager.GenerateTokenPair(person.ID, person.Email)
	if err != nil {
		return nil, fmt.Errorf("generating tokens: %w", err)
	}

	// 4. Create session
	session := &models.Session{
		PersonID:  person.ID,
		TokenHash: s.hashToken(tokens.AccessToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	// Audit Log
	_ = s.auditLogService.Log(ctx, service.LogParams{
		PersonID:     &person.ID,
		Action:       "login",
		ResourceType: "person",
		ResourceID:   person.ID,
		IPAddress:    req.IPAddress,
		UserAgent:    req.UserAgent,
	})

	return &service.LoginResponse{
		User:         person,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int(tokens.ExpiresIn),
	}, nil
}

func (s *authService) Logout(ctx context.Context, token string, ipAddress, userAgent string) error {
	hash := s.hashToken(token)
	session, err := s.authRepo.GetSessionByTokenHash(ctx, hash)
	if err != nil {
		return nil // Already logged out or invalid
	}

	err = s.authRepo.DeleteSession(ctx, session.ID)
	if err == nil {
		_ = s.auditLogService.Log(ctx, service.LogParams{
			PersonID:     &session.PersonID,
			Action:       "logout",
			ResourceType: "person",
			ResourceID:   session.PersonID,
			IPAddress:    ipAddress,
			UserAgent:    userAgent,
		})
	}

	return err
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*service.TokenResponse, error) {
	personID, err := s.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	person, err := s.personRepo.GetByID(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("person not found: %w", err)
	}

	tokens, err := s.tokenManager.GenerateTokenPair(person.ID, person.Email)
	if err != nil {
		return nil, fmt.Errorf("generating tokens: %w", err)
	}

	// Create new session for the new access token
	session := &models.Session{
		PersonID:  person.ID,
		TokenHash: s.hashToken(tokens.AccessToken),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	return &service.TokenResponse{
		AccessToken: tokens.AccessToken,
		ExpiresIn:   int(tokens.ExpiresIn),
	}, nil
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	// Implementation placeholder
	return nil
}

func (s *authService) OAuthLogin(ctx context.Context, provider string, code string) (*service.LoginResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *authService) OAuthCallback(ctx context.Context, provider string, state, code string) (*service.LoginResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *authService) LinkOAuthProvider(ctx context.Context, personID uuid.UUID, provider string, code string) error {
	return errors.New("not implemented")
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	return errors.New("not implemented")
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	return errors.New("not implemented")
}

func (s *authService) ChangePassword(ctx context.Context, personID uuid.UUID, oldPassword, newPassword string) error {
	return errors.New("not implemented")
}

func (s *authService) ValidateSession(ctx context.Context, token string) (*service.SessionInfo, error) {
	claims, err := s.tokenManager.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}

	hash := s.hashToken(token)
	session, err := s.authRepo.GetSessionByTokenHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("session not found or revoked")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		_ = s.auditLogService.Log(ctx, service.LogParams{
			PersonID:     &session.PersonID,
			Action:       "session_expired",
			ResourceType: "person",
			ResourceID:   session.PersonID,
		})
		_ = s.authRepo.DeleteSession(ctx, session.ID)
		return nil, fmt.Errorf("session expired")
	}

	// Update last activity
	session.LastActivity = time.Now()
	_ = s.authRepo.UpdateSession(ctx, session)

	return &service.SessionInfo{
		PersonID:     claims.PersonID,
		Email:        claims.Email,
		ExpiresAt:    session.ExpiresAt,
		LastActivity: session.LastActivity,
	}, nil
}

func (s *authService) GetSessions(ctx context.Context, personID uuid.UUID) ([]*models.Session, error) {
	return s.authRepo.GetSessionsByPerson(ctx, personID)
}

func (s *authService) RevokeSession(ctx context.Context, personID, sessionID uuid.UUID) error {
	return s.authRepo.DeleteSession(ctx, sessionID)
}

func (s *authService) RevokeAllSessions(ctx context.Context, personID uuid.UUID) error {
	return s.authRepo.DeleteSessionsByPerson(ctx, personID)
}

// Helper: Hash token for session storage
func (s *authService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
