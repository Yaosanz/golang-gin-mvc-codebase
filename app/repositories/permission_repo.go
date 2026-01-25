package repositories

import (
	"context"
	"go-starter-app/app/models"

	"github.com/google/uuid"
)

type PermissionRepo struct {
	app IRepoDependencies
}

type IPermissionRepo interface {
	HasPermission(ctx context.Context, roleName string, code string) (bool, error)
	GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error)
}

func NewPermissionRepo(deps IRepoDependencies) IPermissionRepo {
	return &PermissionRepo{
		app: deps,
	}
}

func (r *PermissionRepo) HasPermission(ctx context.Context, roleName string, code string) (bool, error) {
	var count int64
	err := r.app.GetDBWithContext(ctx).
		Table("role_permissions rp").
		Joins("JOIN permissions p ON rp.permission_id = p.id").
		Joins("JOIN roles r ON rp.role_id = r.id").
		Where("r.name = ? AND p.name = ?", roleName, code).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PermissionRepo) GetPermissionsByRoleID(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission

	err := r.app.GetDBWithContext(ctx).
		Table("permissions p").
		Joins("JOIN role_permissions rp ON p.id = rp.permission_id").
		Where("rp.role_id = ?", roleID).
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}
