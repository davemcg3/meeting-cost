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

type personRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewPersonRepository creates a new GORM-based PersonRepository.
func NewPersonRepository(db *gorm.DB, cache cache.Cache) repository.PersonRepository {
	return &personRepository{
		db:    db,
		cache: cache,
	}
}

func (r *personRepository) Create(ctx context.Context, person *models.Person) error {
	if err := r.db.WithContext(ctx).Create(person).Error; err != nil {
		return fmt.Errorf("creating person: %w", err)
	}
	return nil
}

func (r *personRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Person, error) {
	// 1. Check cache
	cacheKey := cache.KeyPerson(id)
	var person models.Person
	if err := r.cache.Get(ctx, cacheKey, &person); err == nil {
		return &person, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).First(&person, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("person not found: %w", err)
		}
		return nil, fmt.Errorf("getting person by id: %w", err)
	}

	// 3. Set cache (TTL: 1 hour for persons as they don't change frequently)
	_ = r.cache.Set(ctx, cacheKey, person, 1*time.Hour)

	return &person, nil
}

func (r *personRepository) GetByEmail(ctx context.Context, email string) (*models.Person, error) {
	// 1. Check cache
	cacheKey := cache.KeyPersonByEmail(email)
	var person models.Person
	if err := r.cache.Get(ctx, cacheKey, &person); err == nil {
		return &person, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).First(&person, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("person not found by email: %w", err)
		}
		return nil, fmt.Errorf("getting person by email: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, person, 1*time.Hour)

	return &person, nil
}

func (r *personRepository) List(ctx context.Context, filters repository.PersonFilters, pagination repository.Pagination) ([]*models.Person, int64, error) {
	var persons []*models.Person
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Person{})

	// Apply filters
	if filters.Email != nil {
		query = query.Where("email = ?", *filters.Email)
	}
	if filters.Anonymized != nil {
		query = query.Where("anonymized = ?", *filters.Anonymized)
	}
	if filters.OrganizationID != nil {
		query = query.Joins("JOIN person_organization_profiles ON person_organization_profiles.person_id = persons.id").
			Where("person_organization_profiles.organization_id = ?", *filters.OrganizationID)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("counting persons: %w", err)
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

	if err := query.Find(&persons).Error; err != nil {
		return nil, 0, fmt.Errorf("querying persons: %w", err)
	}

	return persons, total, nil
}

func (r *personRepository) Update(ctx context.Context, person *models.Person) error {
	if err := r.db.WithContext(ctx).Save(person).Error; err != nil {
		return fmt.Errorf("updating person: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyPerson(person.ID))
	_ = r.cache.Delete(ctx, cache.KeyPersonByEmail(person.Email))

	return nil
}

func (r *personRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// We need to get the person first to know their email for cache invalidation
	person, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&models.Person{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("deleting person: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyPerson(id))
	_ = r.cache.Delete(ctx, cache.KeyPersonByEmail(person.Email))

	return nil
}

func (r *personRepository) Anonymize(ctx context.Context, id uuid.UUID) error {
	person, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now()
	person.FirstName = "Anonymized"
	person.LastName = "Anonymized"
	person.Email = fmt.Sprintf("anonymized-%s@example.com", id.String())
	person.Anonymized = true
	person.AnonymizedAt = &now

	if err := r.Update(ctx, person); err != nil {
		return fmt.Errorf("anonymizing person: %w", err)
	}

	return nil
}

func (r *personRepository) GetOrganizations(ctx context.Context, personID uuid.UUID) ([]*models.Organization, error) {
	var orgs []*models.Organization
	err := r.db.WithContext(ctx).
		Joins("JOIN person_organization_profiles ON person_organization_profiles.organization_id = organizations.id").
		Where("person_organization_profiles.person_id = ?", personID).
		Find(&orgs).Error

	if err != nil {
		return nil, fmt.Errorf("getting organizations for person: %w", err)
	}

	return orgs, nil
}

func (r *personRepository) GetActiveOrganizations(ctx context.Context, personID uuid.UUID) ([]*models.Organization, error) {
	var orgs []*models.Organization
	err := r.db.WithContext(ctx).
		Joins("JOIN person_organization_profiles ON person_organization_profiles.organization_id = organizations.id").
		Where("person_organization_profiles.person_id = ? AND person_organization_profiles.is_active = ?", personID, true).
		Find(&orgs).Error

	if err != nil {
		return nil, fmt.Errorf("getting active organizations for person: %w", err)
	}

	return orgs, nil
}
