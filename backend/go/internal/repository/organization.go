package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
)

// OrganizationRepository handles all database operations for Organization entities.
type OrganizationRepository interface {
	// Create
	Create(ctx context.Context, org *models.Organization) error

	// Read
	GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	GetBySlug(ctx context.Context, slug string) (*models.Organization, error)
	List(ctx context.Context, filters OrgFilters, pagination Pagination) ([]*models.Organization, int64, error)

	// Update
	Update(ctx context.Context, org *models.Organization) error

	// Delete (soft delete)
	Delete(ctx context.Context, id uuid.UUID) error

	// Members
	GetMembers(ctx context.Context, orgID uuid.UUID, activeOnly bool) ([]*models.PersonOrganizationProfile, error)
	AddMember(ctx context.Context, profile *models.PersonOrganizationProfile) error
	RemoveMember(ctx context.Context, personID, orgID uuid.UUID) error
	UpdateMemberProfile(ctx context.Context, profile *models.PersonOrganizationProfile) error

	// Meetings
	GetMeetings(ctx context.Context, orgID uuid.UUID, filters MeetingFilters, pagination Pagination) ([]*models.Meeting, int64, error)
}

type OrgFilters struct {
	Slug     *string
	Name     *string
	MemberID *uuid.UUID // Filter by member
}

