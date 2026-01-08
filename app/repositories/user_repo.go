package repositories

import (
	"context"

	"go-starter-app/app/http/utils"
	"go-starter-app/app/models"
)

// allowedUserFilter defines the fields that can be used for filtering users
var allowedUserFilter = map[string]bool{
	"company_id": true,
	"is_active":  true,
}

type UserRepo struct {
	app IRepoDependencies
}

type IUserRepo interface {
	FindAll(ctx context.Context, params utils.QueryParams) ([]models.User, int64, error)
	FindById(ctx context.Context, id string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	UpdateById(ctx context.Context, id string, updates map[string]interface{}) error
	Delete(ctx context.Context, id string) error
	MarkAsActive(ctx context.Context, id string) (bool, error)
	MarkAsInActive(ctx context.Context, id string) (bool, error)
}

// NewUserRepo creates a new instance of UserRepo
func NewUserRepo(deps IRepoDependencies) IUserRepo {
	return &UserRepo{
		app: deps,
	}
}

func (r *UserRepo) FindAll(ctx context.Context, params utils.QueryParams) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.app.GetDBWithContext(ctx).Model(&models.User{})
	query = query.Scopes(utils.ApplySearch(params.Search, "username", "name", "email", "phone"))
	query = query.Scopes(utils.ApplyFilter(params.Filters, allowedUserFilter))
	query = query.Scopes(utils.ApplyDateFilter(params.DateFrom, params.DateTo, params.DateField))

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Scopes(utils.ApplySorting(params))
	query = query.Scopes(utils.ApplyPagination(params))

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepo) FindById(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.app.GetDBWithContext(ctx).First(&user, "id = ?", id).Error
	return &user, err
}

func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.app.GetDBWithContext(ctx).
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	return r.app.GetDBWithContext(ctx).Create(user).Error
}

func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	return r.app.GetDBWithContext(ctx).Save(user).Error
}

func (r *UserRepo) UpdateById(ctx context.Context, id string, updates map[string]interface{}) error {
	return r.app.GetDBWithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.app.GetDBWithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	return &user, err
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	return r.app.GetDBWithContext(ctx).
		Where("id = ?", id).
		Delete(&models.User{}).Error
}

func (r *UserRepo) MarkAsActive(ctx context.Context, id string) (bool, error) {
	err := r.app.GetDBWithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", true).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *UserRepo) MarkAsInActive(ctx context.Context, id string) (bool, error) {
	err := r.app.GetDBWithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", false).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
