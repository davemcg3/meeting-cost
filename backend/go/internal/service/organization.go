package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// OrganizationService handles organization-related business logic.
type OrganizationService interface {
	// CRUD
	CreateOrganization(ctx context.Context, creatorID uuid.UUID, req CreateOrganizationRequest) (*OrganizationDTO, error)
	GetOrganization(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID) (*OrganizationDTO, error)
	ListOrganizations(ctx context.Context, requesterID uuid.UUID) ([]*OrganizationDTO, error)
	UpdateOrganization(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, req UpdateOrganizationRequest) (*OrganizationDTO, error)
	DeleteOrganization(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, ipAddress, userAgent string) error

	// Members
	GetMembers(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID) ([]*MemberDTO, error)
	AddMember(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, req AddMemberRequest) error
	RemoveMember(ctx context.Context, orgID uuid.UUID, requesterID, memberID uuid.UUID, ipAddress, userAgent string) error
	UpdateMemberWage(ctx context.Context, orgID uuid.UUID, personID uuid.UUID, wage float64, requesterID uuid.UUID, ipAddress, userAgent string) error

	// Settings
	UpdateSettings(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, settings map[string]interface{}) error
	UpdateDefaultWage(ctx context.Context, orgID uuid.UUID, wage float64, requesterID uuid.UUID) error
	SetBlendedWage(ctx context.Context, orgID uuid.UUID, enabled bool, requesterID uuid.UUID) error

	// Permissions
	GetRoles(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID) ([]*RoleDTO, error)
	CreateRole(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, req CreateRoleRequest) (*RoleDTO, error)
	AssignRole(ctx context.Context, orgID uuid.UUID, personID uuid.UUID, roleID uuid.UUID, requesterID uuid.UUID) error
}

type CreateOrganizationRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	DefaultWage float64 `json:"default_wage" validate:"min=0"`
	IPAddress   string  `json:"-"`
	UserAgent   string  `json:"-"`
}

type UpdateOrganizationRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	DefaultWage *float64 `json:"default_wage,omitempty"`
	IPAddress   string   `json:"-"`
	UserAgent   string   `json:"-"`
}

type OrganizationDTO struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Slug           string    `json:"slug"`
	Description    string    `json:"description"`
	DefaultWage    float64   `json:"default_wage"`
	UseBlendedWage bool      `json:"use_blended_wage"`
	CreatedAt      time.Time `json:"created_at"`
	MemberCount    int       `json:"member_count"`
}

type MemberDTO struct {
	PersonID   uuid.UUID `json:"person_id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	IsActive   bool      `json:"is_active"`
	HourlyWage *float64  `json:"hourly_wage,omitempty"` // Only visible to authorized users
	JoinedAt   time.Time `json:"joined_at"`
	Roles      []string  `json:"roles"`
}

type AddMemberRequest struct {
	PersonID  uuid.UUID `json:"person_id"`
	Email     string    `json:"email"`
	Wage      *float64  `json:"wage"`
	IPAddress string    `json:"-"`
	UserAgent string    `json:"-"`
}

type RoleDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateRoleRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"` // e.g., "meeting:create"
}
