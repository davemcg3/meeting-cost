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

type organizationRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewOrganizationRepository creates a new GORM-based OrganizationRepository.
func NewOrganizationRepository(db *gorm.DB, cache cache.Cache) repository.OrganizationRepository {
	return &organizationRepository{
		db:    db,
		cache: cache,
	}
}

func (r *organizationRepository) Create(ctx context.Context, org *models.Organization) error {
	if err := r.db.WithContext(ctx).Create(org).Error; err != nil {
		return fmt.Errorf("creating organization: %w", err)
	}
	return nil
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	// 1. Check cache
	cacheKey := cache.KeyOrganization(id)
	var org models.Organization
	if err := r.cache.Get(ctx, cacheKey, &org); err == nil {
		return &org, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).First(&org, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("organization not found: %w", err)
		}
		return nil, fmt.Errorf("getting organization by id: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, org, 1*time.Hour)

	return &org, nil
}

func (r *organizationRepository) GetBySlug(ctx context.Context, slug string) (*models.Organization, error) {
	// 1. Check cache
	cacheKey := cache.KeyOrganizationBySlug(slug)
	var org models.Organization
	if err := r.cache.Get(ctx, cacheKey, &org); err == nil {
		return &org, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).First(&org, "slug = ?", slug).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("organization not found by slug: %w", err)
		}
		return nil, fmt.Errorf("getting organization by slug: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, org, 1*time.Hour)

	return &org, nil
}

func (r *organizationRepository) List(ctx context.Context, filters repository.OrgFilters, pagination repository.Pagination) ([]*models.Organization, int64, error) {
	var orgs []*models.Organization
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Organization{})

	// Apply filters
	if filters.Name != nil {
		query = query.Where("name ILIKE ?", "%"+*filters.Name+"%")
	}
	if filters.Slug != nil {
		query = query.Where("slug = ?", *filters.Slug)
	}
	if filters.MemberID != nil {
		query = query.Joins("JOIN person_organization_profiles ON person_organization_profiles.organization_id = organizations.id").
			Where("person_organization_profiles.person_id = ?", *filters.MemberID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("counting organizations: %w", err)
	}

	// Apply pagination
	if pagination.PageSize > 0 {
		query = query.Offset(pagination.Offset()).Limit(pagination.Limit())
	}

	// Apply sorting
	if pagination.SortBy != "" {
		sortDir := "ASC"
		if pagination.SortDir == "desc" {
			sortDir = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", pagination.SortBy, sortDir))
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Find(&orgs).Error; err != nil {
		return nil, 0, fmt.Errorf("querying organizations: %w", err)
	}

	return orgs, total, nil
}

func (r *organizationRepository) Update(ctx context.Context, org *models.Organization) error {
	if err := r.db.WithContext(ctx).Save(org).Error; err != nil {
		return fmt.Errorf("updating organization: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyOrganization(org.ID))
	_ = r.cache.Delete(ctx, cache.KeyOrganizationBySlug(org.Slug))

	return nil
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	org, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&models.Organization{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("deleting organization: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyOrganization(id))
	_ = r.cache.Delete(ctx, cache.KeyOrganizationBySlug(org.Slug))

	return nil
}

func (r *organizationRepository) GetMembers(ctx context.Context, orgID uuid.UUID, activeOnly bool) ([]*models.PersonOrganizationProfile, error) {
	var profiles []*models.PersonOrganizationProfile
	query := r.db.WithContext(ctx).Where("organization_id = ?", orgID)
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Preload("Person").Find(&profiles).Error; err != nil {
		return nil, fmt.Errorf("getting organization members: %w", err)
	}

	return profiles, nil
}

func (r *organizationRepository) AddMember(ctx context.Context, profile *models.PersonOrganizationProfile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return fmt.Errorf("adding member to organization: %w", err)
	}

	// Invalidate related cache entries if needed
	// (e.g., if we cached member lists)
	return nil
}

func (r *organizationRepository) RemoveMember(ctx context.Context, personID, orgID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("person_id = ? AND organization_id = ?", personID, orgID).
		Delete(&models.PersonOrganizationProfile{}).Error; err != nil {
		return fmt.Errorf("removing member from organization: %w", err)
	}
	return nil
}

func (r *organizationRepository) UpdateMemberProfile(ctx context.Context, profile *models.PersonOrganizationProfile) error {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return fmt.Errorf("updating member profile: %w", err)
	}
	return nil
}

func (r *organizationRepository) GetMeetings(ctx context.Context, orgID uuid.UUID, filters repository.MeetingFilters, pagination repository.Pagination) ([]*models.Meeting, int64, error) {
	var meetings []*models.Meeting
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Meeting{}).Where("organization_id = ?", orgID)

	// Apply filters
	if filters.CreatedByID != nil {
		query = query.Where("created_by_id = ?", *filters.CreatedByID)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if filters.StartedAfter != nil {
		query = query.Where("started_at >= ?", *filters.StartedAfter)
	}
	if filters.StartedBefore != nil {
		query = query.Where("started_at <= ?", *filters.StartedBefore)
	}
	if filters.ExternalType != nil {
		query = query.Where("external_type = ?", *filters.ExternalType)
	}
	if filters.ExternalID != nil {
		query = query.Where("external_id = ?", *filters.ExternalID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("counting organization meetings: %w", err)
	}

	// Apply pagination
	if pagination.PageSize > 0 {
		query = query.Offset(pagination.Offset()).Limit(pagination.Limit())
	}

	// Apply sorting
	if pagination.SortBy != "" {
		sortDir := "ASC"
		if pagination.SortDir == "desc" {
			sortDir = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", pagination.SortBy, sortDir))
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Find(&meetings).Error; err != nil {
		return nil, 0, fmt.Errorf("querying organization meetings: %w", err)
	}

	return meetings, total, nil
}
