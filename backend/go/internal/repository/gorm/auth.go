package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/cache"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
	"github.com/yourorg/meeting-cost/backend/go/internal/repository"
	"gorm.io/gorm"
)

type authRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewAuthRepository creates a new GORM-based AuthRepository.
func NewAuthRepository(db *gorm.DB, cache cache.Cache) repository.AuthRepository {
	return &authRepository{
		db:    db,
		cache: cache,
	}
}

// AuthMethod operations

func (r *authRepository) CreateAuthMethod(ctx context.Context, method *models.AuthMethod) error {
	if err := r.db.WithContext(ctx).Create(method).Error; err != nil {
		return fmt.Errorf("creating auth method: %w", err)
	}
	return nil
}

func (r *authRepository) GetAuthMethodByID(ctx context.Context, id uuid.UUID) (*models.AuthMethod, error) {
	// 1. Check cache
	cacheKey := cache.KeyAuthMethod(id)
	var method models.AuthMethod
	if err := r.cache.Get(ctx, cacheKey, &method); err == nil {
		return &method, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).First(&method, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("auth method not found: %w", err)
		}
		return nil, fmt.Errorf("getting auth method by id: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, method, 1*time.Hour)

	return &method, nil
}

func (r *authRepository) GetAuthMethodByProvider(ctx context.Context, provider, providerID string) (*models.AuthMethod, error) {
	// 1. Check cache
	cacheKey := cache.KeyAuthMethodByProvider(provider, providerID)
	var method models.AuthMethod
	if err := r.cache.Get(ctx, cacheKey, &method); err == nil {
		return &method, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).First(&method, "provider = ? AND provider_id = ?", provider, providerID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("auth method not found by provider: %w", err)
		}
		return nil, fmt.Errorf("getting auth method by provider: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, method, 1*time.Hour)

	return &method, nil
}

func (r *authRepository) GetAuthMethodsByPerson(ctx context.Context, personID uuid.UUID) ([]*models.AuthMethod, error) {
	var methods []*models.AuthMethod
	if err := r.db.WithContext(ctx).Where("person_id = ?", personID).Find(&methods).Error; err != nil {
		return nil, fmt.Errorf("getting auth methods by person: %w", err)
	}
	return methods, nil
}

func (r *authRepository) UpdateAuthMethod(ctx context.Context, method *models.AuthMethod) error {
	if err := r.db.WithContext(ctx).Save(method).Error; err != nil {
		return fmt.Errorf("updating auth method: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyAuthMethod(method.ID))
	_ = r.cache.Delete(ctx, cache.KeyAuthMethodByProvider(method.Provider, method.ProviderID))

	return nil
}

func (r *authRepository) DeleteAuthMethod(ctx context.Context, id uuid.UUID) error {
	method, err := r.GetAuthMethodByID(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&models.AuthMethod{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("deleting auth method: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyAuthMethod(id))
	_ = r.cache.Delete(ctx, cache.KeyAuthMethodByProvider(method.Provider, method.ProviderID))

	return nil
}

// Session operations

func (r *authRepository) CreateSession(ctx context.Context, session *models.Session) error {
	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return fmt.Errorf("creating session: %w", err)
	}
	return nil
}

func (r *authRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*models.Session, error) {
	// 1. Check cache
	cacheKey := cache.KeySession(tokenHash)
	var session models.Session
	if err := r.cache.Get(ctx, cacheKey, &session); err == nil {
		return &session, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).First(&session, "token_hash = ?", tokenHash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("session not found: %w", err)
		}
		return nil, fmt.Errorf("getting session by token hash: %w", err)
	}

	// 3. Set cache
	// Calculate remaining TTL
	remaining := time.Until(session.ExpiresAt)
	if remaining > 0 {
		_ = r.cache.Set(ctx, cacheKey, session, remaining)
	}

	return &session, nil
}

func (r *authRepository) GetSessionsByPerson(ctx context.Context, personID uuid.UUID) ([]*models.Session, error) {
	var sessions []*models.Session
	if err := r.db.WithContext(ctx).Where("person_id = ?", personID).Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("getting sessions by person: %w", err)
	}
	return sessions, nil
}

func (r *authRepository) UpdateSession(ctx context.Context, session *models.Session) error {
	if err := r.db.WithContext(ctx).Save(session).Error; err != nil {
		return fmt.Errorf("updating session: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeySession(session.TokenHash))

	return nil
}

func (r *authRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	var session models.Session
	if err := r.db.WithContext(ctx).First(&session, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // Already deleted
		}
		return fmt.Errorf("getting session for deletion: %w", err)
	}

	if err := r.db.WithContext(ctx).Delete(&models.Session{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeySession(session.TokenHash))

	return nil
}

func (r *authRepository) DeleteExpiredSessions(ctx context.Context) error {
	// Not ideal for cache invalidation as we don't know the hashes,
	// but expired sessions shouldn't be in cache due to TTL.
	if err := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error; err != nil {
		return fmt.Errorf("deleting expired sessions: %w", err)
	}
	return nil
}

func (r *authRepository) DeleteSessionsByPerson(ctx context.Context, personID uuid.UUID) error {
	var sessions []*models.Session
	if err := r.db.WithContext(ctx).Where("person_id = ?", personID).Find(&sessions).Error; err != nil {
		return fmt.Errorf("getting sessions for bulk deletion: %w", err)
	}

	if err := r.db.WithContext(ctx).Where("person_id = ?", personID).Delete(&models.Session{}).Error; err != nil {
		return fmt.Errorf("deleting sessions by person: %w", err)
	}

	// Invalidate cache for each session
	for _, s := range sessions {
		_ = r.cache.Delete(ctx, cache.KeySession(s.TokenHash))
	}

	return nil
}
