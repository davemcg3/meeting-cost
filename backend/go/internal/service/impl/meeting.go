package impl

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/cache"
	"github.com/yourorg/meeting-cost/backend/go/internal/logger"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
	"github.com/yourorg/meeting-cost/backend/go/internal/pubsub"
	"github.com/yourorg/meeting-cost/backend/go/internal/repository"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
)

type meetingService struct {
	meetingRepo     repository.MeetingRepository
	incrementRepo   repository.IncrementRepository
	orgRepo         repository.OrganizationRepository
	profileRepo     repository.PersonOrganizationProfileRepository
	permissionRepo  repository.PermissionRepository
	auditLogService service.AuditLogService
	cache           cache.Cache
	pubsub          pubsub.PubSub
	logger          logger.Logger
}

// NewMeetingService creates a new MeetingService implementation.
func NewMeetingService(
	meetingRepo repository.MeetingRepository,
	incrementRepo repository.IncrementRepository,
	orgRepo repository.OrganizationRepository,
	profileRepo repository.PersonOrganizationProfileRepository,
	permissionRepo repository.PermissionRepository,
	auditLogService service.AuditLogService,
	cache cache.Cache,
	ps pubsub.PubSub,
	logger logger.Logger,
) service.MeetingService {
	return &meetingService{
		meetingRepo:     meetingRepo,
		incrementRepo:   incrementRepo,
		orgRepo:         orgRepo,
		profileRepo:     profileRepo,
		permissionRepo:  permissionRepo,
		auditLogService: auditLogService,
		cache:           cache,
		pubsub:          ps,
		logger:          logger,
	}
}

func (s *meetingService) broadcastEvent(ctx context.Context, meetingID uuid.UUID, eventType service.EventType, payload interface{}) {
	event := service.MeetingEvent{
		Type:      eventType,
		MeetingID: meetingID,
		Payload:   payload,
	}

	channel := cache.ChannelMeetingEvents(meetingID)
	if err := s.pubsub.Publish(ctx, channel, event); err != nil {
		s.logger.Error("failed to broadcast meeting event", "meeting_id", meetingID, "type", eventType, "error", err)
	}
}

func (s *meetingService) CreateMeeting(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, req service.CreateMeetingRequest) (*service.MeetingDTO, error) {
	// 1. Authorization check
	hasPermission, err := s.permissionRepo.HasPermission(ctx, requesterID, orgID, "meeting", nil, "create")
	if err != nil {
		return nil, fmt.Errorf("checking permission: %w", err)
	}
	if !hasPermission {
		return nil, fmt.Errorf("forbidden: insufficient permissions to create meeting")
	}

	// 2. Business validation (e.g. org exists and is active)
	if _, err := s.orgRepo.GetByID(ctx, orgID); err != nil {
		return nil, fmt.Errorf("getting organization: %w", err)
	}

	// 3. Create model
	meeting := &models.Meeting{
		OrganizationID: orgID,
		CreatedByID:    requesterID,
		Purpose:        req.Purpose,
		ExternalType:   req.ExternalType,
		ExternalID:     req.ExternalID,
		IsActive:       false,
	}

	// 4. Repository call
	if err := s.meetingRepo.Create(ctx, meeting); err != nil {
		return nil, fmt.Errorf("creating meeting: %w", err)
	}

	// Audit Log
	_ = s.auditLogService.Log(ctx, service.LogParams{
		PersonID:       &requesterID,
		OrganizationID: &meeting.OrganizationID,
		Action:         "create_meeting",
		ResourceType:   "meeting",
		ResourceID:     meeting.ID,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
	})

	// 5. Return DTO
	return s.toMeetingDTO(meeting), nil
}

func (s *meetingService) GetMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID) (*service.MeetingDTO, error) {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return nil, err
	}

	// Authorization check
	hasPermission, err := s.permissionRepo.HasPermission(ctx, requesterID, meeting.OrganizationID, "meeting", &meetingID, "read")
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, fmt.Errorf("forbidden")
	}

	return s.toMeetingDTO(meeting), nil
}

func (s *meetingService) UpdateMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID, req service.UpdateMeetingRequest) (*service.MeetingDTO, error) {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return nil, err
	}

	// Authorization check
	hasPermission, err := s.permissionRepo.HasPermission(ctx, requesterID, meeting.OrganizationID, "meeting", &meetingID, "update")
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, fmt.Errorf("forbidden")
	}

	if req.Purpose != nil {
		meeting.Purpose = *req.Purpose
	}

	if err := s.meetingRepo.Update(ctx, meeting); err != nil {
		return nil, err
	}

	return s.toMeetingDTO(meeting), nil
}

func (s *meetingService) DeleteMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID, ipAddress, userAgent string) error {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return err
	}

	// Authorization check
	hasPermission, err := s.permissionRepo.HasPermission(ctx, requesterID, meeting.OrganizationID, "meeting", &meetingID, "delete")
	if err != nil {
		return err
	}
	if !hasPermission {
		return fmt.Errorf("forbidden")
	}

	err = s.meetingRepo.Delete(ctx, meetingID)
	if err == nil {
		_ = s.auditLogService.Log(ctx, service.LogParams{
			PersonID:       &requesterID,
			OrganizationID: &meeting.OrganizationID,
			Action:         "delete_meeting",
			ResourceType:   "meeting",
			ResourceID:     meetingID,
			IPAddress:      ipAddress,
			UserAgent:      userAgent,
		})
	}
	return err
}

func (s *meetingService) StartMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID) error {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return err
	}

	// Authorization check
	hasPermission, err := s.permissionRepo.HasPermission(ctx, requesterID, meeting.OrganizationID, "meeting", &meetingID, "start")
	if err != nil {
		return err
	}
	if !hasPermission {
		return fmt.Errorf("forbidden")
	}

	if meeting.IsActive {
		return fmt.Errorf("meeting is already active")
	}

	if err := s.meetingRepo.Start(ctx, meetingID); err != nil {
		return err
	}

	// Create first increment
	org, _ := s.orgRepo.GetByID(ctx, meeting.OrganizationID)
	firstInc := &models.Increment{
		MeetingID:     meetingID,
		StartTime:     time.Now(),
		AverageWage:   org.DefaultWage,
		AttendeeCount: 0, // Should probably be based on current participants if any
		Purpose:       meeting.Purpose,
	}

	if err := s.meetingRepo.AddIncrement(ctx, firstInc); err != nil {
		return err
	}

	s.broadcastEvent(ctx, meetingID, service.EventMeetingStarted, firstInc)
	return nil
}

func (s *meetingService) StopMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID) error {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return err
	}

	// Authorization check
	hasPermission, err := s.permissionRepo.HasPermission(ctx, requesterID, meeting.OrganizationID, "meeting", &meetingID, "stop")
	if err != nil {
		return err
	}
	if !hasPermission {
		return fmt.Errorf("forbidden")
	}

	if !meeting.IsActive {
		return fmt.Errorf("meeting is not active")
	}

	if err := s.meetingRepo.Stop(ctx, meetingID); err != nil {
		return err
	}

	// Finalize current increment
	increments, _ := s.meetingRepo.GetIncrements(ctx, meetingID)
	now := time.Now()
	for _, inc := range increments {
		if inc.StopTime.IsZero() {
			inc.StopTime = now
			inc.ElapsedTime = int(now.Sub(inc.StartTime).Seconds())
			inc.Cost = (float64(inc.ElapsedTime) / 3600.0) * float64(inc.AttendeeCount) * inc.AverageWage
			_ = s.incrementRepo.Update(ctx, inc)
			break
		}
	}

	// Update meeting totals
	if err := s.updateMeetingTotals(ctx, meetingID); err != nil {
		s.logger.Error("failed to update meeting totals on stop", "meeting_id", meetingID, "error", err)
	}

	s.broadcastEvent(ctx, meetingID, service.EventMeetingStopped, nil)
	return nil
}

func (s *meetingService) ResetMeeting(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID) error {
	// Implementation for resetting a meeting
	return nil
}

func (s *meetingService) UpdateAttendeeCount(ctx context.Context, meetingID uuid.UUID, count int, requesterID uuid.UUID, ipAddress, userAgent string) error {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return err
	}

	// Auth check
	hasPerm, _ := s.permissionRepo.HasPermission(ctx, requesterID, meeting.OrganizationID, "meeting", &meetingID, "update")
	if !hasPerm {
		return fmt.Errorf("forbidden")
	}

	if !meeting.IsActive {
		// Just update the meeting record if not active
		return nil // Or update a default count if we add it
	}

	err = s.cycleIncrement(ctx, meetingID, func(inc *models.Increment) {
		inc.AttendeeCount = count
	})

	if err == nil {
		_ = s.auditLogService.Log(ctx, service.LogParams{
			PersonID:       &requesterID,
			OrganizationID: &meeting.OrganizationID,
			Action:         "update_attendee_count",
			ResourceType:   "meeting",
			ResourceID:     meetingID,
			Details:        map[string]interface{}{"attendee_count": count},
			IPAddress:      ipAddress,
			UserAgent:      userAgent,
		})
	}

	return err
}

func (s *meetingService) UpdateAverageWage(ctx context.Context, meetingID uuid.UUID, wage float64, requesterID uuid.UUID) error {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return err
	}

	hasPerm, _ := s.permissionRepo.HasPermission(ctx, requesterID, meeting.OrganizationID, "meeting", &meetingID, "update")
	if !hasPerm {
		return fmt.Errorf("forbidden")
	}

	if !meeting.IsActive {
		return nil
	}

	return s.cycleIncrement(ctx, meetingID, func(inc *models.Increment) {
		inc.AverageWage = wage
	})
}

func (s *meetingService) UpdatePurpose(ctx context.Context, meetingID uuid.UUID, purpose string, requesterID uuid.UUID) error {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return err
	}

	hasPerm, _ := s.permissionRepo.HasPermission(ctx, requesterID, meeting.OrganizationID, "meeting", &meetingID, "update")
	if !hasPerm {
		return fmt.Errorf("forbidden")
	}

	if !meeting.IsActive {
		meeting.Purpose = purpose
		return s.meetingRepo.Update(ctx, meeting)
	}

	return s.cycleIncrement(ctx, meetingID, func(inc *models.Increment) {
		inc.Purpose = purpose
	})
}

// cycleIncrement stops the current increment and starts a new one with modifications
func (s *meetingService) cycleIncrement(ctx context.Context, meetingID uuid.UUID, modify func(*models.Increment)) error {
	increments, err := s.meetingRepo.GetIncrements(ctx, meetingID)
	if err != nil {
		return err
	}

	now := time.Now()
	var lastInc *models.Increment
	for _, inc := range increments {
		if inc.StopTime.IsZero() {
			lastInc = inc
			break
		}
	}

	newInc := &models.Increment{
		MeetingID: meetingID,
		StartTime: now,
	}

	if lastInc != nil {
		lastInc.StopTime = now
		lastInc.ElapsedTime = int(now.Sub(lastInc.StartTime).Seconds())
		// Basic cost calculation: (elapsed / 3600) * count * average_wage
		lastInc.Cost = (float64(lastInc.ElapsedTime) / 3600.0) * float64(lastInc.AttendeeCount) * lastInc.AverageWage

		if err := s.incrementRepo.Update(ctx, lastInc); err != nil {
			return err
		}

		// Inherit values from last increment
		newInc.AttendeeCount = lastInc.AttendeeCount
		newInc.AverageWage = lastInc.AverageWage
		newInc.Purpose = lastInc.Purpose
	} else {
		// No active increment? Fallback to meeting defaults or current state
		meeting, _ := s.meetingRepo.GetByID(ctx, meetingID)
		org, _ := s.orgRepo.GetByID(ctx, meeting.OrganizationID)
		newInc.AverageWage = org.DefaultWage
		newInc.Purpose = meeting.Purpose
	}

	modify(newInc)

	if err := s.meetingRepo.AddIncrement(ctx, newInc); err != nil {
		return err
	}

	// Update meeting totals
	if err := s.updateMeetingTotals(ctx, meetingID); err != nil {
		s.logger.Error("failed to update meeting totals on cycle", "meeting_id", meetingID, "error", err)
	}

	s.broadcastEvent(ctx, meetingID, service.EventMeetingCost, newInc)
	return nil
}

func (s *meetingService) AddParticipant(ctx context.Context, meetingID uuid.UUID, personID uuid.UUID, requesterID uuid.UUID) error {
	// Implementation for adding participant
	return nil
}

func (s *meetingService) RemoveParticipant(ctx context.Context, meetingID uuid.UUID, personID uuid.UUID, requesterID uuid.UUID) error {
	// Implementation for removing participant
	return nil
}

func (s *meetingService) ListMeetings(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, filters service.MeetingFilters, pagination service.Pagination) ([]*service.MeetingDTO, int64, error) {
	// Authorization check: must be a member of the organization
	profile, err := s.profileRepo.GetByPersonAndOrg(ctx, requesterID, orgID)
	if err != nil || !profile.IsActive {
		return nil, 0, fmt.Errorf("forbidden: not a member of this organization")
	}

	repoFilters := repository.MeetingFilters{
		OrganizationID: &orgID,
		IsActive:       filters.IsActive,
		StartedAfter:   filters.StartedAfter,
		StartedBefore:  filters.StartedBefore,
	}

	repoPagination := repository.Pagination{
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
	}

	meetings, total, err := s.meetingRepo.List(ctx, repoFilters, repoPagination)
	if err != nil {
		return nil, 0, fmt.Errorf("listing meetings: %w", err)
	}

	dtos := make([]*service.MeetingDTO, len(meetings))
	for i, m := range meetings {
		dtos[i] = s.toMeetingDTO(m)
	}

	return dtos, total, nil
}

func (s *meetingService) GetMeetingCost(ctx context.Context, meetingID uuid.UUID, requesterID uuid.UUID) (*service.MeetingCostDTO, error) {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return nil, err
	}

	increments, err := s.meetingRepo.GetIncrements(ctx, meetingID)
	if err != nil {
		return nil, err
	}

	var totalCost float64
	var totalDuration int
	now := time.Now()

	for _, inc := range increments {
		if !inc.StopTime.IsZero() {
			totalCost += inc.Cost
			totalDuration += inc.ElapsedTime
		} else if meeting.IsActive {
			// Current active increment
			elapsed := int(now.Sub(inc.StartTime).Seconds())
			currentCost := (float64(elapsed) / 3600.0) * float64(inc.AttendeeCount) * inc.AverageWage
			totalCost += currentCost
			totalDuration += elapsed
		}
	}

	res := &service.MeetingCostDTO{
		TotalCost:     totalCost,
		TotalDuration: totalDuration,
	}

	if totalDuration > 0 {
		res.CostPerSecond = totalCost / float64(totalDuration)
		res.CostPerMinute = res.CostPerSecond * 60
		res.CostPerHour = res.CostPerSecond * 3600
	}

	return res, nil
}

func (s *meetingService) DeduplicateMeeting(ctx context.Context, meetingID uuid.UUID, externalType, externalID string) (*service.MeetingDTO, error) {
	// Implementation for deduplicating meeting
	return nil, nil
}

// Helper methods

// toMeetingDTO converts a meeting model to a DTO.
func (s *meetingService) toMeetingDTO(m *models.Meeting) *service.MeetingDTO {
	return &service.MeetingDTO{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Purpose:        m.Purpose,
		StartedAt:      m.StartedAt,
		StoppedAt:      m.StoppedAt,
		IsActive:       m.IsActive,
		TotalCost:      m.TotalCost,
		TotalDuration:  m.TotalDuration,
		MaxAttendees:   m.MaxAttendees,
		CreatedAt:      m.CreatedAt,
	}
}

// updateMeetingTotals recalculates and updates the meeting's cached total fields.
func (s *meetingService) updateMeetingTotals(ctx context.Context, meetingID uuid.UUID) error {
	meeting, err := s.meetingRepo.GetByID(ctx, meetingID)
	if err != nil {
		return fmt.Errorf("getting meeting: %w", err)
	}

	increments, err := s.meetingRepo.GetIncrements(ctx, meetingID)
	if err != nil {
		return fmt.Errorf("getting increments: %w", err)
	}

	var totalCost float64
	var totalDuration int
	var maxAttendees int

	for _, inc := range increments {
		if !inc.StopTime.IsZero() {
			totalCost += inc.Cost
			totalDuration += inc.ElapsedTime
			if inc.AttendeeCount > maxAttendees {
				maxAttendees = inc.AttendeeCount
			}

			// Update increment's running total cost if needed
			if inc.TotalCost != totalCost {
				inc.TotalCost = totalCost
				_ = s.incrementRepo.Update(ctx, inc)
			}
		} else if meeting.IsActive {
			// Current active increment - it should also contribute to max attendees
			if inc.AttendeeCount > maxAttendees {
				maxAttendees = inc.AttendeeCount
			}
		}
	}

	meeting.TotalCost = totalCost
	meeting.TotalDuration = totalDuration
	meeting.MaxAttendees = maxAttendees

	if err := s.meetingRepo.Update(ctx, meeting); err != nil {
		return fmt.Errorf("updating meeting totals: %w", err)
	}

	return nil
}
