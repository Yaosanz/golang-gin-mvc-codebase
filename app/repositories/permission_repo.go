package repositories

import (
	"context"

	"github.com/google/uuid"
)

type PermissionRepo struct {
	app IRepoDependencies
}

type IPermissionRepo interface {
	HasPermission(ctx context.Context, roleID uuid.UUID, code string) (bool, error)
}

func NewPermissionRepo(deps IRepoDependencies) IPermissionRepo {
	return &PermissionRepo{
		app: deps,
	}
}

func (r *PermissionRepo) HasPermission(ctx context.Context, roleID uuid.UUID, code string) (bool, error) {
	var count int64
	err := r.app.GetDBWithContext(ctx).
		Table("role_permissions rp").
		Joins("JOIN permissions p ON rp.permission_id = p.id").
		Where("rp.role_id = ? AND p.name = ?", roleID, code).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
