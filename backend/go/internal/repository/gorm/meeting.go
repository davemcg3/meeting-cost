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

type meetingRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewMeetingRepository creates a new GORM-based MeetingRepository.
func NewMeetingRepository(db *gorm.DB, cache cache.Cache) repository.MeetingRepository {
	return &meetingRepository{
		db:    db,
		cache: cache,
	}
}

func (r *meetingRepository) Create(ctx context.Context, meeting *models.Meeting) error {
	if err := r.db.WithContext(ctx).Create(meeting).Error; err != nil {
		return fmt.Errorf("creating meeting: %w", err)
	}
	return nil
}

func (r *meetingRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Meeting, error) {
	// 1. Check cache
	cacheKey := cache.KeyMeeting(id)
	var meeting models.Meeting
	if err := r.cache.Get(ctx, cacheKey, &meeting); err == nil {
		return &meeting, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).First(&meeting, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("meeting not found: %w", err)
		}
		return nil, fmt.Errorf("getting meeting by id: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, meeting, 15*time.Minute)

	return &meeting, nil
}

func (r *meetingRepository) GetByExternalID(ctx context.Context, externalType, externalID string) (*models.Meeting, error) {
	// 1. Check cache
	cacheKey := cache.KeyMeetingByExternalID(externalType, externalID)
	var meeting models.Meeting
	if err := r.cache.Get(ctx, cacheKey, &meeting); err == nil {
		return &meeting, nil
	}

	// 2. Query DB
	if err := r.db.WithContext(ctx).First(&meeting, "external_type = ? AND external_id = ?", externalType, externalID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("meeting not found by external id: %w", err)
		}
		return nil, fmt.Errorf("getting meeting by external id: %w", err)
	}

	// 3. Set cache
	_ = r.cache.Set(ctx, cacheKey, meeting, 15*time.Minute)

	return &meeting, nil
}

func (r *meetingRepository) GetByDeduplicationHash(ctx context.Context, hash string) (*models.Meeting, error) {
	var meeting models.Meeting
	if err := r.db.WithContext(ctx).First(&meeting, "deduplication_hash = ?", hash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("meeting not found by deduplication hash: %w", err)
		}
		return nil, fmt.Errorf("getting meeting by deduplication hash: %w", err)
	}
	return &meeting, nil
}

func (r *meetingRepository) List(ctx context.Context, filters repository.MeetingFilters, pagination repository.Pagination) ([]*models.Meeting, int64, error) {
	var meetings []*models.Meeting
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Meeting{})

	// Apply filters
	if filters.OrganizationID != nil {
		query = query.Where("organization_id = ?", *filters.OrganizationID)
	}
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
		return nil, 0, fmt.Errorf("counting meetings: %w", err)
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
		return nil, 0, fmt.Errorf("querying meetings: %w", err)
	}

	return meetings, total, nil
}

func (r *meetingRepository) Update(ctx context.Context, meeting *models.Meeting) error {
	if err := r.db.WithContext(ctx).Save(meeting).Error; err != nil {
		return fmt.Errorf("updating meeting: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyMeeting(meeting.ID))
	if meeting.ExternalID != "" {
		_ = r.cache.Delete(ctx, cache.KeyMeetingByExternalID(meeting.ExternalType, meeting.ExternalID))
	}

	return nil
}

func (r *meetingRepository) Start(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&models.Meeting{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active":  true,
			"started_at": &now,
		}).Error

	if err != nil {
		return fmt.Errorf("starting meeting: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyMeeting(id))
	// External ID cache would also need invalidation if we want to be thorough
	return nil
}

func (r *meetingRepository) Stop(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	err := r.db.WithContext(ctx).Model(&models.Meeting{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active": false,
			"stopped_at":  &now,
		}).Error

	if err != nil {
		return fmt.Errorf("stopping meeting: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyMeeting(id))
	return nil
}

func (r *meetingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	meeting, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&models.Meeting{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("deleting meeting: %w", err)
	}

	// Invalidate cache
	_ = r.cache.Delete(ctx, cache.KeyMeeting(id))
	if meeting.ExternalID != "" {
		_ = r.cache.Delete(ctx, cache.KeyMeetingByExternalID(meeting.ExternalType, meeting.ExternalID))
	}

	return nil
}

func (r *meetingRepository) GetIncrements(ctx context.Context, meetingID uuid.UUID) ([]*models.Increment, error) {
	var increments []*models.Increment
	if err := r.db.WithContext(ctx).Where("meeting_id = ?", meetingID).Order("start_time ASC").Find(&increments).Error; err != nil {
		return nil, fmt.Errorf("getting increments: %w", err)
	}
	return increments, nil
}

func (r *meetingRepository) AddIncrement(ctx context.Context, increment *models.Increment) error {
	if err := r.db.WithContext(ctx).Create(increment).Error; err != nil {
		return fmt.Errorf("adding increment: %w", err)
	}
	return nil
}

func (r *meetingRepository) GetParticipants(ctx context.Context, meetingID uuid.UUID) ([]*models.MeetingParticipant, error) {
	var participants []*models.MeetingParticipant
	if err := r.db.WithContext(ctx).Where("meeting_id = ?", meetingID).Preload("Person").Find(&participants).Error; err != nil {
		return nil, fmt.Errorf("getting participants: %w", err)
	}
	return participants, nil
}

func (r *meetingRepository) AddParticipant(ctx context.Context, participant *models.MeetingParticipant) error {
	if err := r.db.WithContext(ctx).Create(participant).Error; err != nil {
		return fmt.Errorf("adding participant: %w", err)
	}
	return nil
}

func (r *meetingRepository) RemoveParticipant(ctx context.Context, meetingID, personID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("meeting_id = ? AND person_id = ?", meetingID, personID).
		Delete(&models.MeetingParticipant{}).Error; err != nil {
		return fmt.Errorf("removing participant: %w", err)
	}
	return nil
}
