package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// MeetingService handles meeting-related business logic.
type MeetingService interface {
	// CRUD
	CreateMeeting(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, req CreateMeetingRequest) (*MeetingDTO, error)
	GetMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID) (*MeetingDTO, error)
	UpdateMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID, req UpdateMeetingRequest) (*MeetingDTO, error)
	DeleteMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID, ipAddress, userAgent string) error

	// Meeting control
	StartMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID) error
	StopMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID) error
	ResetMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID) error

	// Increments
	UpdateAttendeeCount(ctx context.Context, meetingID uuid.UUID, count int, requesterID uuid.UUID, ipAddress, userAgent string) error
	UpdateAverageWage(ctx context.Context, meetingID uuid.UUID, wage float64, requesterID uuid.UUID) error
	UpdatePurpose(ctx context.Context, meetingID uuid.UUID, purpose string, requesterID uuid.UUID) error

	// Participants
	AddParticipant(ctx context.Context, meetingID uuid.UUID, personID uuid.UUID, requesterID uuid.UUID) error
	RemoveParticipant(ctx context.Context, meetingID uuid.UUID, personID uuid.UUID, requesterID uuid.UUID) error

	// Queries
	ListMeetings(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, filters MeetingFilters, pagination Pagination) ([]*MeetingDTO, int64, error)
	GetMeetingCost(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID) (*MeetingCostDTO, error)

	// Deduplication
	DeduplicateMeeting(ctx context.Context, meetingID uuid.UUID, externalType, externalID string) (*MeetingDTO, error)
}

type CreateMeetingRequest struct {
	OrganizationID uuid.UUID `json:"organization_id" validate:"required"`
	Purpose        string    `json:"purpose"`
	ExternalType   string    `json:"external_type"` // "zoom", "teams", etc.
	ExternalID     string    `json:"external_id"`
	IPAddress      string    `json:"-"`
	UserAgent      string    `json:"-"`
}

type UpdateMeetingRequest struct {
	Purpose *string `json:"purpose"`
}

type MeetingDTO struct {
	ID             uuid.UUID        `json:"id"`
	OrganizationID uuid.UUID        `json:"organization_id"`
	Purpose        string           `json:"purpose"`
	StartedAt      *time.Time       `json:"started_at"`
	StoppedAt      *time.Time       `json:"stopped_at"`
	IsActive       bool             `json:"is_active"`
	TotalCost      float64          `json:"total_cost"`
	TotalDuration  int              `json:"total_duration"` // seconds
	MaxAttendees   int              `json:"max_attendees"`
	Increments     []IncrementDTO   `json:"increments,omitempty"`
	Participants   []ParticipantDTO `json:"participants,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
}

type IncrementDTO struct {
	ID            uuid.UUID `json:"id"`
	StartTime     time.Time `json:"start_time"`
	StopTime      time.Time `json:"stop_time"`
	ElapsedTime   int       `json:"elapsed_time"` // seconds
	AttendeeCount int       `json:"attendee_count"`
	AverageWage   float64   `json:"average_wage"`
	Cost          float64   `json:"cost"`
	TotalCost     float64   `json:"total_cost"`
	Purpose       string    `json:"purpose"`
}

type ParticipantDTO struct {
	PersonID uuid.UUID  `json:"person_id"`
	Email    string     `json:"email"`
	Name     string     `json:"name"`
	JoinedAt *time.Time `json:"joined_at"`
	LeftAt   *time.Time `json:"left_at"`
}

type MeetingCostDTO struct {
	TotalCost     float64 `json:"total_cost"`
	TotalDuration int     `json:"total_duration"` // seconds
	CostPerSecond float64 `json:"cost_per_second"`
	CostPerMinute float64 `json:"cost_per_minute"`
	CostPerHour   float64 `json:"cost_per_hour"`
}

// MeetingFilters here mirrors repository.MeetingFilters, but is kept separate
// so the service API remains decoupled from repository concerns.
type MeetingFilters struct {
	IsActive      *bool
	StartedAfter  *time.Time
	StartedBefore *time.Time
}

// Pagination is reused from the repository layer for convenience.
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}
