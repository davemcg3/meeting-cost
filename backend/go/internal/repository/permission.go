package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/yourorg/meeting-cost/backend/go/internal/models"
)

// PermissionRepository handles permission and role operations.
type PermissionRepository interface {
	// Role operations
	CreateRole(ctx context.Context, role *models.Role) error
	GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	GetRolesByOrganization(ctx context.Context, orgID uuid.UUID) ([]*models.Role, error)
	UpdateRole(ctx context.Context, role *models.Role) error
	DeleteRole(ctx context.Context, id uuid.UUID) error

	// Permission operations
	CreatePermission(ctx context.Context, permission *models.Permission) error
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	GetPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*models.Permission, error)
	GetPermissionsByPerson(ctx context.Context, personID uuid.UUID) ([]*models.Permission, error)
	GetPermissionsByOrganization(ctx context.Context, orgID uuid.UUID) ([]*models.Permission, error)
	UpdatePermission(ctx context.Context, permission *models.Permission) error
	DeletePermission(ctx context.Context, id uuid.UUID) error

	// Role assignment
	AssignRole(ctx context.Context, assignment *models.RoleAssignment) error
	UnassignRole(ctx context.Context, roleID, personID, orgID uuid.UUID) error
	GetRolesByPerson(ctx context.Context, personID, orgID uuid.UUID) ([]*models.Role, error)

	// Permission checking
	HasPermission(ctx context.Context, personID, orgID uuid.UUID, resourceName string, resourceID *uuid.UUID, activity string) (bool, error)
}

