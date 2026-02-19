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

type incrementRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewIncrementRepository creates a new GORM-based IncrementRepository.
func NewIncrementRepository(db *gorm.DB, cache cache.Cache) repository.IncrementRepository {
	return &incrementRepository{
		db:    db,
		cache: cache,
	}
}

func (r *incrementRepository) Create(ctx context.Context, increment *models.Increment) error {
	if err := r.db.WithContext(ctx).Create(increment).Error; err != nil {
		return fmt.Errorf("creating increment: %w", err)
	}
	// Invalidate increments list for meeting
	_ = r.cache.Delete(ctx, cache.KeyMeetingIncrements(increment.MeetingID))
	return nil
}

func (r *incrementRepository) CreateBatch(ctx context.Context, increments []*models.Increment) error {
	if len(increments) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).Create(&increments).Error; err != nil {
		return fmt.Errorf("batch creating increments: %w", err)
	}
	// Invalidate increments list for meeting
	_ = r.cache.Delete(ctx, cache.KeyMeetingIncrements(increments[0].MeetingID))
	return nil
}

func (r *incrementRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Increment, error) {
	// 1. Check cache
	cacheKey := cache.KeyIncrement(id)
	var increment models.Increment
	if err := r.cache.Get(ctx, cacheKey, &increment); err == nil {
		return &increment, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).First(&increment, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("increment not found: %w", err)
		}
		return nil, fmt.Errorf("getting increment by id: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, increment, 1*time.Hour)

	return &increment, nil
}

func (r *incrementRepository) GetByMeeting(ctx context.Context, meetingID uuid.UUID) ([]*models.Increment, error) {
	// 1. Check cache
	cacheKey := cache.KeyMeetingIncrements(meetingID)
	var increments []*models.Increment
	if err := r.cache.Get(ctx, cacheKey, &increments); err == nil {
		return increments, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).Where("meeting_id = ?", meetingID).Order("start_time ASC").Find(&increments).Error; err != nil {
		return nil, fmt.Errorf("getting increments by meeting: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, increments, 15*time.Minute)

	return increments, nil
}

func (r *incrementRepository) Update(ctx context.Context, increment *models.Increment) error {
	if err := r.db.WithContext(ctx).Save(increment).Error; err != nil {
		return fmt.Errorf("updating increment: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyIncrement(increment.ID))
	_ = r.cache.Delete(ctx, cache.KeyMeetingIncrements(increment.MeetingID))

	return nil
}

func (r *incrementRepository) Delete(ctx context.Context, id uuid.UUID) error {
	inc, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&models.Increment{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("deleting increment: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyIncrement(id))
	_ = r.cache.Delete(ctx, cache.KeyMeetingIncrements(inc.MeetingID))

	return nil
}

func (r *incrementRepository) DeleteByMeeting(ctx context.Context, meetingID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("meeting_id = ?", meetingID).Delete(&models.Increment{}).Error; err != nil {
		return fmt.Errorf("deleting increments by meeting: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyMeetingIncrements(meetingID))
	// Individual increments are still in cache, but harder to invalidate without list.
	// In a real scenario, we might want to iterate and delete or use a different strategy.
	return nil
}
