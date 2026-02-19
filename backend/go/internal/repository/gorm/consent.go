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

type consentRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewConsentRepository creates a new GORM-based ConsentRepository.
func NewConsentRepository(db *gorm.DB, cache cache.Cache) repository.ConsentRepository {
	return &consentRepository{
		db:    db,
		cache: cache,
	}
}

func (r *consentRepository) Create(ctx context.Context, consent *models.CookieConsent) error {
	if err := r.db.WithContext(ctx).Create(consent).Error; err != nil {
		return fmt.Errorf("creating consent: %w", err)
	}
	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyConsentBySession(consent.SessionID))
	if consent.PersonID != nil {
		_ = r.cache.Delete(ctx, cache.KeyConsentByPerson(*consent.PersonID))
	}
	return nil
}

func (r *consentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.CookieConsent, error) {
	var consent models.CookieConsent
	if err := r.db.WithContext(ctx).First(&consent, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("consent not found: %w", err)
		}
		return nil, fmt.Errorf("getting consent by id: %w", err)
	}
	return &consent, nil
}

func (r *consentRepository) GetCurrentBySession(ctx context.Context, sessionID string) (*models.CookieConsent, error) {
	// 1. Check cache
	cacheKey := cache.KeyConsentBySession(sessionID)
	var consent models.CookieConsent
	if err := r.cache.Get(ctx, cacheKey, &consent); err == nil {
		return &consent, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).Order("created_at DESC").First(&consent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("consent not found for session: %w", err)
		}
		return nil, fmt.Errorf("getting current consent by session: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, consent, 1*time.Hour)

	return &consent, nil
}

func (r *consentRepository) GetCurrentByPerson(ctx context.Context, personID uuid.UUID) (*models.CookieConsent, error) {
	// 1. Check cache
	cacheKey := cache.KeyConsentByPerson(personID)
	var consent models.CookieConsent
	if err := r.cache.Get(ctx, cacheKey, &consent); err == nil {
		return &consent, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).Where("person_id = ?", personID).Order("created_at DESC").First(&consent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("consent not found for person: %w", err)
		}
		return nil, fmt.Errorf("getting current consent by person: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, consent, 1*time.Hour)

	return &consent, nil
}

func (r *consentRepository) GetHistoryBySession(ctx context.Context, sessionID string) ([]*models.CookieConsent, error) {
	var history []*models.CookieConsent
	if err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).Order("created_at DESC").Find(&history).Error; err != nil {
		return nil, fmt.Errorf("getting consent history by session: %w", err)
	}
	return history, nil
}

func (r *consentRepository) GetHistoryByPerson(ctx context.Context, personID uuid.UUID) ([]*models.CookieConsent, error) {
	var history []*models.CookieConsent
	if err := r.db.WithContext(ctx).Where("person_id = ?", personID).Order("created_at DESC").Find(&history).Error; err != nil {
		return nil, fmt.Errorf("getting consent history by person: %w", err)
	}
	return history, nil
}

func (r *consentRepository) Update(ctx context.Context, consent *models.CookieConsent) error {
	if err := r.db.WithContext(ctx).Save(consent).Error; err != nil {
		return fmt.Errorf("updating consent: %w", err)
	}
	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyConsentBySession(consent.SessionID))
	if consent.PersonID != nil {
		_ = r.cache.Delete(ctx, cache.KeyConsentByPerson(*consent.PersonID))
	}
	return nil
}

func (r *consentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	consent, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&models.CookieConsent{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("deleting consent: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyConsentBySession(consent.SessionID))
	if consent.PersonID != nil {
		_ = r.cache.Delete(ctx, cache.KeyConsentByPerson(*consent.PersonID))
	}
	return nil
}
