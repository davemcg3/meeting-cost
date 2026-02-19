package impl

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/logger"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
	"github.com/yourorg/meeting-cost/backend/go/internal/repository"
	"github.com/yourorg/meeting-cost/backend/go/internal/service"
)

type organizationService struct {
	orgRepo         repository.OrganizationRepository
	profileRepo     repository.PersonOrganizationProfileRepository
	permissionRepo  repository.PermissionRepository
	personRepo      repository.PersonRepository
	auditLogService service.AuditLogService
	logger          logger.Logger
}

// NewOrganizationService creates a new OrganizationService implementation.
func NewOrganizationService(
	orgRepo repository.OrganizationRepository,
	profileRepo repository.PersonOrganizationProfileRepository,
	permissionRepo repository.PermissionRepository,
	personRepo repository.PersonRepository,
	auditLogService service.AuditLogService,
	logger logger.Logger,
) service.OrganizationService {
	return &organizationService{
		orgRepo:         orgRepo,
		profileRepo:     profileRepo,
		permissionRepo:  permissionRepo,
		personRepo:      personRepo,
		auditLogService: auditLogService,
		logger:          logger,
	}
}

func (s *organizationService) CreateOrganization(ctx context.Context, creatorID uuid.UUID, req service.CreateOrganizationRequest) (*service.OrganizationDTO, error) {
	// 1. Create model
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	org := &models.Organization{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		DefaultWage: req.DefaultWage,
	}

	// 2. Repository call
	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, fmt.Errorf("creating organization: %w", err)
	}

	// 3. Create initial membership for creator
	profile := &models.PersonOrganizationProfile{
		PersonID:       creatorID,
		OrganizationID: org.ID,
		IsActive:       true,
		HourlyWage:     &req.DefaultWage,
		JoinedAt:       org.CreatedAt,
	}
	if err := s.profileRepo.Create(ctx, profile); err != nil {
		return nil, fmt.Errorf("creating creator profile: %w", err)
	}

	// 4. Seed default roles and assign Admin to creator
	adminRole, err := s.seedDefaultRoles(ctx, org.ID)
	if err != nil {
		s.logger.Error("failed to seed default roles", "org_id", org.ID, "error", err)
	} else if adminRole != nil {
		err = s.permissionRepo.AssignRole(ctx, &models.RoleAssignment{
			RoleID:         adminRole.ID,
			PersonID:       creatorID,
			OrganizationID: org.ID,
		})
		if err != nil {
			s.logger.Error("failed to assign admin role", "org_id", org.ID, "person_id", creatorID, "error", err)
		}
	}

	// Audit Log
	_ = s.auditLogService.Log(ctx, service.LogParams{
		PersonID:       &creatorID,
		OrganizationID: &org.ID,
		Action:         "create_organization",
		ResourceType:   "organization",
		ResourceID:     org.ID,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
	})

	return s.toOrganizationDTO(ctx, org), nil
}

func (s *organizationService) seedDefaultRoles(ctx context.Context, orgID uuid.UUID) (*models.Role, error) {
	// 1. Create Admin Role
	adminRole := &models.Role{
		Name:           "Admin",
		Description:    "Full access to the organization",
		OrganizationID: orgID,
	}
	if err := s.permissionRepo.CreateRole(ctx, adminRole); err != nil {
		return nil, err
	}

	// 2. Create Member Role
	memberRole := &models.Role{
		Name:           "Member",
		Description:    "Standard access to meetings",
		OrganizationID: orgID,
	}
	if err := s.permissionRepo.CreateRole(ctx, memberRole); err != nil {
		return adminRole, err
	}

	// 3. Define Permissions
	perms := []struct {
		RoleID   uuid.UUID
		Resource string
		Activity string
	}{
		// Admin permissions
		{adminRole.ID, "organization", "read"},
		{adminRole.ID, "organization", "update"},
		{adminRole.ID, "organization", "manage_members"},
		{adminRole.ID, "organization", "delete"},
		{adminRole.ID, "meeting", "create"},
		{adminRole.ID, "meeting", "read"},
		{adminRole.ID, "meeting", "update"},
		{adminRole.ID, "meeting", "delete"},
		{adminRole.ID, "meeting", "start"},
		{adminRole.ID, "meeting", "stop"},

		// Member permissions
		{memberRole.ID, "organization", "read"},
		{memberRole.ID, "meeting", "create"},
		{memberRole.ID, "meeting", "read"},
		{memberRole.ID, "meeting", "update"}, // Can update their own meetings (checked in logic)
		{memberRole.ID, "meeting", "start"},
		{memberRole.ID, "meeting", "stop"},
	}

	for _, p := range perms {
		perm := &models.Permission{
			ResourceType:   "role",
			ResourceID:     p.RoleID,
			ResourceName:   p.Resource,
			Activity:       p.Activity,
			Allowed:        true,
			OrganizationID: orgID,
		}
		if err := s.permissionRepo.CreatePermission(ctx, perm); err != nil {
			s.logger.Error("failed to create permission", "role_id", p.RoleID, "resource", p.Resource, "activity", p.Activity, "error", err)
		}
	}

	return adminRole, nil
}

func (s *organizationService) GetOrganization(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID) (*service.OrganizationDTO, error) {
	// Authorization check: requester must be a member
	profile, err := s.profileRepo.GetByPersonAndOrg(ctx, requesterID, orgID)
	if err != nil || !profile.IsActive {
		return nil, fmt.Errorf("forbidden: not a member of this organization")
	}

	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return s.toOrganizationDTO(ctx, org), nil
}

func (s *organizationService) ListOrganizations(ctx context.Context, requesterID uuid.UUID) ([]*service.OrganizationDTO, error) {
	// Filter by member ID
	filters := repository.OrgFilters{
		MemberID: &requesterID,
	}

	// Default pagination (get all for now or first 100)
	pagination := repository.Pagination{Page: 1, PageSize: 100}

	orgs, _, err := s.orgRepo.List(ctx, filters, pagination)
	if err != nil {
		return nil, fmt.Errorf("listing organizations: %w", err)
	}

	dtos := make([]*service.OrganizationDTO, len(orgs))
	for i, org := range orgs {
		dtos[i] = s.toOrganizationDTO(ctx, org)
	}

	return dtos, nil
}

func (s *organizationService) UpdateOrganization(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, req service.UpdateOrganizationRequest) (*service.OrganizationDTO, error) {
	// Authorization check: must have 'update' permission
	hasPerm, err := s.permissionRepo.HasPermission(ctx, requesterID, orgID, "organization", nil, "update")
	if err != nil || !hasPerm {
		return nil, fmt.Errorf("forbidden")
	}

	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Description != nil {
		org.Description = *req.Description
	}
	if req.DefaultWage != nil {
		org.DefaultWage = *req.DefaultWage
	}

	if err := s.orgRepo.Update(ctx, org); err != nil {
		return nil, err
	}

	// Audit Log
	_ = s.auditLogService.Log(ctx, service.LogParams{
		PersonID:       &requesterID,
		OrganizationID: &orgID,
		Action:         "update_organization",
		ResourceType:   "organization",
		ResourceID:     orgID,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
	})

	return s.toOrganizationDTO(ctx, org), nil
}

func (s *organizationService) DeleteOrganization(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, ipAddress, userAgent string) error {
	hasPerm, err := s.permissionRepo.HasPermission(ctx, requesterID, orgID, "organization", nil, "delete")
	if err != nil || !hasPerm {
		return fmt.Errorf("forbidden")
	}

	err = s.orgRepo.Delete(ctx, orgID)
	if err == nil {
		_ = s.auditLogService.Log(ctx, service.LogParams{
			PersonID:       &requesterID,
			OrganizationID: &orgID,
			Action:         "delete_organization",
			ResourceType:   "organization",
			ResourceID:     orgID,
			IPAddress:      ipAddress,
			UserAgent:      userAgent,
		})
	}
	return err
}

func (s *organizationService) GetMembers(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID) ([]*service.MemberDTO, error) {
	// 1. Authorization check: requester must be a member
	profile, err := s.profileRepo.GetByPersonAndOrg(ctx, requesterID, orgID)
	if err != nil || !profile.IsActive {
		return nil, fmt.Errorf("forbidden: not a member of this organization")
	}

	// 2. Fetch all profiles for the org
	profiles, err := s.profileRepo.GetByOrganization(ctx, orgID, false)
	if err != nil {
		return nil, fmt.Errorf("fetching profiles: %w", err)
	}

	// 3. Map to DTOs
	members := make([]*service.MemberDTO, len(profiles))
	for i, p := range profiles {
		members[i] = &service.MemberDTO{
			PersonID:  p.PersonID,
			Email:     p.Person.Email,
			FirstName: p.Person.FirstName,
			LastName:  p.Person.LastName,
			IsActive:  p.IsActive,
			JoinedAt:  p.JoinedAt,
		}

		// Auth check for wage visibility (admin vs self)
		if requesterID == p.PersonID {
			members[i].HourlyWage = p.HourlyWage
		} else {
			// Check if requester is admin
			isAdmin, _ := s.permissionRepo.HasPermission(ctx, requesterID, orgID, "organization", nil, "manage_members")
			if isAdmin {
				members[i].HourlyWage = p.HourlyWage
			}
		}
	}

	return members, nil
}

func (s *organizationService) AddMember(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, req service.AddMemberRequest) error {
	// 1. Authorization check: must have 'manage_members' permission
	hasPerm, err := s.permissionRepo.HasPermission(ctx, requesterID, orgID, "organization", nil, "manage_members")
	if err != nil || !hasPerm {
		return fmt.Errorf("forbidden")
	}

	// 2. Check if person exists
	var person *models.Person
	if req.PersonID != uuid.Nil {
		person, err = s.personRepo.GetByID(ctx, req.PersonID)
	} else if req.Email != "" {
		person, err = s.personRepo.GetByEmail(ctx, req.Email)
	} else {
		return fmt.Errorf("either person_id or email is required")
	}

	if err != nil {
		return fmt.Errorf("person not found")
	}
	req.PersonID = person.ID

	// 3. Check if already a member
	existing, _ := s.profileRepo.GetByPersonAndOrg(ctx, req.PersonID, orgID)
	if existing != nil {
		if existing.IsActive {
			return fmt.Errorf("person is already a member")
		}
		// Reactivate
		return s.profileRepo.Activate(ctx, req.PersonID, orgID)
	}

	// 4. Create profile
	org, err := s.orgRepo.GetByID(ctx, orgID)
	if err != nil {
		return fmt.Errorf("org not found")
	}

	wage := org.DefaultWage
	if req.Wage != nil {
		wage = *req.Wage
	}

	profile := &models.PersonOrganizationProfile{
		PersonID:       req.PersonID,
		OrganizationID: orgID,
		IsActive:       true,
		HourlyWage:     &wage,
	}

	err = s.profileRepo.Create(ctx, profile)
	if err != nil {
		return err
	}

	// 5. Assign default Member role
	roles, _ := s.permissionRepo.GetRolesByOrganization(ctx, orgID)
	var memberRoleID *uuid.UUID
	for _, r := range roles {
		if r.Name == "Member" {
			memberRoleID = &r.ID
			break
		}
	}

	if memberRoleID != nil {
		_ = s.permissionRepo.AssignRole(ctx, &models.RoleAssignment{
			RoleID:         *memberRoleID,
			PersonID:       req.PersonID,
			OrganizationID: orgID,
		})
	}

	// Audit Log
	_ = s.auditLogService.Log(ctx, service.LogParams{
		PersonID:       &requesterID,
		OrganizationID: &orgID,
		Action:         "add_member",
		ResourceType:   "person",
		ResourceID:     req.PersonID,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
	})
	return nil
}

func (s *organizationService) RemoveMember(ctx context.Context, orgID uuid.UUID, requesterID, memberID uuid.UUID, ipAddress, userAgent string) error {
	// Authorization: must have 'manage_members' or be self
	if requesterID != memberID {
		hasPerm, err := s.permissionRepo.HasPermission(ctx, requesterID, orgID, "organization", nil, "manage_members")
		if err != nil || !hasPerm {
			return fmt.Errorf("forbidden")
		}
	}

	err := s.profileRepo.Deactivate(ctx, memberID, orgID)
	if err == nil {
		_ = s.auditLogService.Log(ctx, service.LogParams{
			PersonID:       &requesterID,
			OrganizationID: &orgID,
			Action:         "remove_member",
			ResourceType:   "person",
			ResourceID:     memberID,
			IPAddress:      ipAddress,
			UserAgent:      userAgent,
		})
	}
	return err
}

func (s *organizationService) UpdateMemberWage(ctx context.Context, orgID uuid.UUID, personID uuid.UUID, wage float64, requesterID uuid.UUID, ipAddress, userAgent string) error {
	// Authorization: must have 'manage_members'
	hasPerm, err := s.permissionRepo.HasPermission(ctx, requesterID, orgID, "organization", nil, "manage_members")
	if err != nil || !hasPerm {
		return fmt.Errorf("forbidden")
	}

	err = s.profileRepo.UpdateWage(ctx, personID, orgID, wage)
	if err == nil {
		_ = s.auditLogService.Log(ctx, service.LogParams{
			PersonID:       &requesterID,
			OrganizationID: &orgID,
			Action:         "update_member_wage",
			ResourceType:   "person",
			ResourceID:     personID,
			Details:        map[string]interface{}{"wage": wage},
			IPAddress:      ipAddress,
			UserAgent:      userAgent,
		})
	}
	return err
}

func (s *organizationService) UpdateSettings(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, settings map[string]interface{}) error {
	return nil
}

func (s *organizationService) UpdateDefaultWage(ctx context.Context, orgID uuid.UUID, wage float64, requesterID uuid.UUID) error {
	return nil
}

func (s *organizationService) SetBlendedWage(ctx context.Context, orgID uuid.UUID, enabled bool, requesterID uuid.UUID) error {
	return nil
}

func (s *organizationService) GetRoles(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID) ([]*service.RoleDTO, error) {
	return nil, nil
}

func (s *organizationService) CreateRole(ctx context.Context, orgID uuid.UUID, requesterID uuid.UUID, req service.CreateRoleRequest) (*service.RoleDTO, error) {
	return nil, nil
}

func (s *organizationService) AssignRole(ctx context.Context, orgID uuid.UUID, personID uuid.UUID, roleID uuid.UUID, requesterID uuid.UUID) error {
	return nil
}

func (s *organizationService) toOrganizationDTO(ctx context.Context, org *models.Organization) *service.OrganizationDTO {
	dto := &service.OrganizationDTO{
		ID:             org.ID,
		Name:           org.Name,
		Slug:           org.Slug,
		Description:    org.Description,
		DefaultWage:    org.DefaultWage,
		UseBlendedWage: org.UseBlendedWage,
		CreatedAt:      org.CreatedAt,
	}

	// Fetch active member count
	// Note: In a high-traffic app, we'd cache this or store it in the org table
	profiles, err := s.profileRepo.GetByOrganization(ctx, org.ID, true)
	if err == nil {
		dto.MemberCount = len(profiles)
	}

	return dto
}
