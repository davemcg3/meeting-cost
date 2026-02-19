package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// PersonService handles person-related business logic.
type PersonService interface {
	// CRUD
	GetPerson(ctx context.Context, personID uuid.UUID) (*PersonDTO, error)
	UpdatePerson(ctx context.Context, personID uuid.UUID, req UpdatePersonRequest) (*PersonDTO, error)

	// Profile
	GetProfile(ctx context.Context, personID uuid.UUID) (*PersonProfileDTO, error)
	UpdateProfile(ctx context.Context, personID uuid.UUID, req UpdateProfileRequest) (*PersonProfileDTO, error)

	// Organizations
	GetOrganizations(ctx context.Context, personID uuid.UUID) ([]*OrganizationDTO, error)
	JoinOrganization(ctx context.Context, personID uuid.UUID, orgID uuid.UUID) error
	LeaveOrganization(ctx context.Context, personID uuid.UUID, orgID uuid.UUID) error

	// GDPR
	RequestDataExport(ctx context.Context, personID uuid.UUID) (*DataExportResponse, error)
	RequestDeletion(ctx context.Context, personID uuid.UUID) error

	// Settings
	UpdateSettings(ctx context.Context, personID uuid.UUID, settings map[string]interface{}) error
}

type PersonDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}

type PersonProfileDTO struct {
	PersonDTO
	Organizations []OrganizationMembershipDTO `json:"organizations"`
	AuthMethods   []AuthMethodDTO             `json:"auth_methods"`
}

type AuthMethodDTO struct {
	ID                 uuid.UUID `json:"id"`
	Provider           string    `json:"provider"`
	ProviderIdentifier string    `json:"provider_identifier"`
	CreatedAt          time.Time `json:"created_at"`
}

type OrganizationMembershipDTO struct {
	OrganizationID   uuid.UUID `json:"organization_id"`
	OrganizationName string    `json:"organization_name"`
	IsActive         bool      `json:"is_active"`
	JoinedAt         time.Time `json:"joined_at"`
	Role             string    `json:"role,omitempty"`
}

type UpdatePersonRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Timezone  *string `json:"timezone"`
	Locale    *string `json:"locale"`
}

// UpdateProfileRequest and DataExportResponse are placeholders; full shape will
// be defined when implementing PersonService.
type UpdateProfileRequest struct {
	// Extend with profile-specific fields as needed.
}

type DataExportResponse struct {
	PersonID uuid.UUID `json:"person_id"`
	// Add exported data fields/URLs when implementing export.
}

