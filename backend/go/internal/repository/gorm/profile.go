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

type profileRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewPersonOrganizationProfileRepository creates a new GORM-based PersonOrganizationProfileRepository.
func NewPersonOrganizationProfileRepository(db *gorm.DB, cache cache.Cache) repository.PersonOrganizationProfileRepository {
	return &profileRepository{
		db:    db,
		cache: cache,
	}
}

func (r *profileRepository) Create(ctx context.Context, profile *models.PersonOrganizationProfile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return fmt.Errorf("creating profile: %w", err)
	}
	return nil
}

func (r *profileRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.PersonOrganizationProfile, error) {
	// 1. Check cache
	cacheKey := cache.KeyProfile(id)
	var profile models.PersonOrganizationProfile
	if err := r.cache.Get(ctx, cacheKey, &profile); err == nil {
		return &profile, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).First(&profile, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("profile not found: %w", err)
		}
		return nil, fmt.Errorf("getting profile by id: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, profile, 1*time.Hour)

	return &profile, nil
}

func (r *profileRepository) GetByPersonAndOrg(ctx context.Context, personID, orgID uuid.UUID) (*models.PersonOrganizationProfile, error) {
	// 1. Check cache
	cacheKey := cache.KeyProfileByPersonAndOrg(personID, orgID)
	var profile models.PersonOrganizationProfile
	if err := r.cache.Get(ctx, cacheKey, &profile); err == nil {
		return &profile, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).Where("person_id = ? AND organization_id = ?", personID, orgID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("profile not found by person and org: %w", err)
		}
		return nil, fmt.Errorf("getting profile by person and org: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, profile, 1*time.Hour)

	return &profile, nil
}

func (r *profileRepository) GetByPerson(ctx context.Context, personID uuid.UUID) ([]*models.PersonOrganizationProfile, error) {
	var profiles []*models.PersonOrganizationProfile
	if err := r.db.WithContext(ctx).Where("person_id = ?", personID).Find(&profiles).Error; err != nil {
		return nil, fmt.Errorf("getting profiles by person: %w", err)
	}
	return profiles, nil
}

func (r *profileRepository) GetByOrganization(ctx context.Context, orgID uuid.UUID, activeOnly bool) ([]*models.PersonOrganizationProfile, error) {
	var profiles []*models.PersonOrganizationProfile
	query := r.db.WithContext(ctx).Preload("Person").Where("organization_id = ?", orgID)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}
	if err := query.Find(&profiles).Error; err != nil {
		return nil, fmt.Errorf("getting profiles by organization: %w", err)
	}
	return profiles, nil
}

func (r *profileRepository) Update(ctx context.Context, profile *models.PersonOrganizationProfile) error {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return fmt.Errorf("updating profile: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyProfile(profile.ID))
	_ = r.cache.Delete(ctx, cache.KeyProfileByPersonAndOrg(profile.PersonID, profile.OrganizationID))

	return nil
}

func (r *profileRepository) UpdateWage(ctx context.Context, personID, orgID uuid.UUID, wage float64) error {
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&models.PersonOrganizationProfile{}).
		Where("person_id = ? AND organization_id = ?", personID, orgID).
		Updates(map[string]interface{}{
			"hourly_wage":     wage,
			"wage_updated_at": &now,
		}).Error

	if err != nil {
		return fmt.Errorf("updating wage: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyProfileByPersonAndOrg(personID, orgID))

	return nil
}

func (r *profileRepository) Activate(ctx context.Context, personID, orgID uuid.UUID) error {
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&models.PersonOrganizationProfile{}).
		Where("person_id = ? AND organization_id = ?", personID, orgID).
		Updates(map[string]interface{}{
			"is_active": true,
			"left_at":   nil,
			"joined_at": now, // Reset joined at or just reactivate? Let's follow business logic.
		}).Error

	if err != nil {
		return fmt.Errorf("activating profile: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyProfileByPersonAndOrg(personID, orgID))

	return nil
}

func (r *profileRepository) Deactivate(ctx context.Context, personID, orgID uuid.UUID) error {
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&models.PersonOrganizationProfile{}).
		Where("person_id = ? AND organization_id = ?", personID, orgID).
		Updates(map[string]interface{}{
			"is_active": false,
			"left_at":   &now,
		}).Error

	if err != nil {
		return fmt.Errorf("deactivating profile: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyProfileByPersonAndOrg(personID, orgID))

	return nil
}

func (r *profileRepository) Delete(ctx context.Context, id uuid.UUID) error {
	profile, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&models.PersonOrganizationProfile{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("deleting profile: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyProfile(id))
	_ = r.cache.Delete(ctx, cache.KeyProfileByPersonAndOrg(profile.PersonID, profile.OrganizationID))

	return nil
}
