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

type permissionRepository struct {
	db    *gorm.DB
	cache cache.Cache
}

// NewPermissionRepository creates a new GORM-based PermissionRepository.
func NewPermissionRepository(db *gorm.DB, cache cache.Cache) repository.PermissionRepository {
	return &permissionRepository{
		db:    db,
		cache: cache,
	}
}

// Role operations

func (r *permissionRepository) CreateRole(ctx context.Context, role *models.Role) error {
	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		return fmt.Errorf("creating role: %w", err)
	}
	return nil
}

func (r *permissionRepository) GetRoleByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	cacheKey := cache.KeyRole(id)
	var role models.Role
	if err := r.cache.Get(ctx, cacheKey, &role); err == nil {
		return &role, nil
	}

	if err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("role not found: %w", err)
		}
		return nil, fmt.Errorf("getting role by id: %w", err)
	}

	_ = r.cache.Set(ctx, cacheKey, role, 1*time.Hour)
	return &role, nil
}

func (r *permissionRepository) GetRolesByOrganization(ctx context.Context, orgID uuid.UUID) ([]*models.Role, error) {
	var roles []*models.Role
	// Roles can be global (orgID is null) or org-specific
	if err := r.db.WithContext(ctx).Where("organization_id = ? OR organization_id IS NULL", orgID).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("getting roles by organization: %w", err)
	}
	return roles, nil
}

func (r *permissionRepository) UpdateRole(ctx context.Context, role *models.Role) error {
	if err := r.db.WithContext(ctx).Save(role).Error; err != nil {
		return fmt.Errorf("updating role: %w", err)
	}
	_ = r.cache.Delete(ctx, cache.KeyRole(role.ID))
	return nil
}

func (r *permissionRepository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Role{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("deleting role: %w", err)
	}
	_ = r.cache.Delete(ctx, cache.KeyRole(id))
	return nil
}

// Permission operations

func (r *permissionRepository) CreatePermission(ctx context.Context, permission *models.Permission) error {
	if err := r.db.WithContext(ctx).Create(permission).Error; err != nil {
		return fmt.Errorf("creating permission: %w", err)
	}
	return nil
}

func (r *permissionRepository) GetPermissionByID(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	cacheKey := cache.KeyPermission(id)
	var permission models.Permission
	if err := r.cache.Get(ctx, cacheKey, &permission); err == nil {
		return &permission, nil
	}

	if err := r.db.WithContext(ctx).First(&permission, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("getting permission by id: %w", err)
	}

	_ = r.cache.Set(ctx, cacheKey, permission, 1*time.Hour)
	return &permission, nil
}

func (r *permissionRepository) GetPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*models.Permission, error) {
	var permissions []*models.Permission
	if err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("getting permissions by role: %w", err)
	}
	return permissions, nil
}

func (r *permissionRepository) GetPermissionsByPerson(ctx context.Context, personID uuid.UUID) ([]*models.Permission, error) {
	var permissions []*models.Permission
	// This would get permissions from all roles assigned to the person across all organizations
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN role_assignments ON role_assignments.role_id = role_permissions.role_id").
		Where("role_assignments.person_id = ?", personID).
		Find(&permissions).Error

	if err != nil {
		return nil, fmt.Errorf("getting permissions by person: %w", err)
	}
	return permissions, nil
}

func (r *permissionRepository) GetPermissionsByOrganization(ctx context.Context, orgID uuid.UUID) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN roles ON roles.id = role_permissions.role_id").
		Where("roles.organization_id = ? OR roles.organization_id IS NULL", orgID).
		Find(&permissions).Error

	if err != nil {
		return nil, fmt.Errorf("getting permissions by organization: %w", err)
	}
	return permissions, nil
}

func (r *permissionRepository) UpdatePermission(ctx context.Context, permission *models.Permission) error {
	if err := r.db.WithContext(ctx).Save(permission).Error; err != nil {
		return fmt.Errorf("updating permission: %w", err)
	}
	_ = r.cache.Delete(ctx, cache.KeyPermission(permission.ID))
	return nil
}

func (r *permissionRepository) DeletePermission(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.Permission{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("deleting permission: %w", err)
	}
	_ = r.cache.Delete(ctx, cache.KeyPermission(id))
	return nil
}

// Role assignment

func (r *permissionRepository) AssignRole(ctx context.Context, assignment *models.RoleAssignment) error {
	if err := r.db.WithContext(ctx).Create(assignment).Error; err != nil {
		return fmt.Errorf("assigning role: %w", err)
	}
	// Invalidate permission checks for this user
	// (Hard to invalidate specific ones without knowing resource/activity)
	return nil
}

func (r *permissionRepository) UnassignRole(ctx context.Context, roleID, personID, orgID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("role_id = ? AND person_id = ? AND (organization_id = ? OR organization_id IS NULL)", roleID, personID, orgID).
		Delete(&models.RoleAssignment{}).Error; err != nil {
		return fmt.Errorf("unassigning role: %w", err)
	}
	return nil
}

func (r *permissionRepository) GetRolesByPerson(ctx context.Context, personID, orgID uuid.UUID) ([]*models.Role, error) {
	var roles []*models.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN role_assignments ON role_assignments.role_id = roles.id").
		Where("role_assignments.person_id = ? AND (role_assignments.organization_id = ? OR role_assignments.organization_id IS NULL)", personID, orgID).
		Find(&roles).Error

	if err != nil {
		return nil, fmt.Errorf("getting roles by person: %w", err)
	}
	return roles, nil
}

// Permission checking

func (r *permissionRepository) HasPermission(ctx context.Context, personID, orgID uuid.UUID, resourceName string, resourceID *uuid.UUID, activity string) (bool, error) {
	// 1. Check cache
	cacheKey := cache.KeyHasPermission(personID, orgID, resourceName, resourceID, activity)
	var hasPermission bool
	if err := r.cache.Get(ctx, cacheKey, &hasPermission); err == nil {
		return hasPermission, nil
	}

	// 2. Query DB
	// We check if any role assigned to the person in this org has the required permission,
	// OR if the person has the permission directly assigned.
	var count int64

	// Query for role-based permissions
	roleQuery := r.db.WithContext(ctx).
		Table("permissions").
		Select("count(*)").
		Joins("JOIN role_assignments ON role_assignments.role_id = permissions.resource_id").
		Where("permissions.resource_type = ?", "role").
		Where("role_assignments.person_id = ?", personID).
		Where("(role_assignments.organization_id = ? OR role_assignments.organization_id IS NULL)", orgID).
		Where("permissions.resource_name = ? AND permissions.activity = ?", resourceName, activity).
		Where("permissions.allowed = ?", true)

	if resourceID != nil {
		roleQuery = roleQuery.Where("(permissions.target_resource_id = ? OR permissions.target_resource_id IS NULL)", *resourceID)
	} else {
		roleQuery = roleQuery.Where("permissions.target_resource_id IS NULL")
	}

	if err := roleQuery.Count(&count).Error; err != nil {
		return false, fmt.Errorf("checking role-based permission: %w", err)
	}

	if count > 0 {
		hasPermission = true
	} else {
		// Query for person-direct permissions
		personQuery := r.db.WithContext(ctx).
			Table("permissions").
			Where("resource_type = ? AND resource_id = ?", "person", personID).
			Where("(organization_id = ? OR organization_id IS NULL)", orgID).
			Where("resource_name = ? AND activity = ?", resourceName, activity).
			Where("allowed = ?", true)

		if resourceID != nil {
			personQuery = personQuery.Where("(target_resource_id = ? OR target_resource_id IS NULL)", *resourceID)
		} else {
			personQuery = personQuery.Where("target_resource_id IS NULL")
		}

		if err := personQuery.Count(&count).Error; err != nil {
			return false, fmt.Errorf("checking person-based permission: %w", err)
		}
		hasPermission = count > 0
	}

	// 3. Set cache (Short TTL as permissions might change)
	_ = r.cache.Set(ctx, cacheKey, hasPermission, 1*time.Minute)

	return hasPermission, nil
}
