package repositories

import (
	"context"

	"go-starter-app/app/models"

	"github.com/google/uuid"
)

type RoleRepo struct {
	app IRepoDependencies
}

type IRoleRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error)
	GetByName(ctx context.Context, name string) (*models.Role, error)
	GetAll(ctx context.Context) ([]models.Role, error)
	Create(ctx context.Context, role *models.Role) error
	Update(ctx context.Context, role *models.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewRoleRepo(deps IRepoDependencies) IRoleRepo {
	return &RoleRepo{
		app: deps,
	}
}

func (r *RoleRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.app.GetDBWithContext(ctx).First(&role, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) GetByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.app.GetDBWithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) GetAll(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	err := r.app.GetDBWithContext(ctx).Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) Create(ctx context.Context, role *models.Role) error {
	return r.app.GetDBWithContext(ctx).Create(role).Error
}

func (r *RoleRepo) Update(ctx context.Context, role *models.Role) error {
	return r.app.GetDBWithContext(ctx).Save(role).Error
}

func (r *RoleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.app.GetDBWithContext(ctx).Delete(&models.Role{}, "id = ?", id).Error
}
