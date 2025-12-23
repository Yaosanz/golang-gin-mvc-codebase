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
	FindById(ctx context.Context, id int64) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	UpdateById(ctx context.Context, id int64, updates map[string]interface{}) error
	Delete(ctx context.Context, id int64) error
	MarkAsActive(ctx context.Context, id int64) (bool, error)
	MarkAsInActive(ctx context.Context, id int64) (bool, error)
}

// NewUserRepo creates a new instance of UserRepo
func NewUserRepo(deps IRepoDependencies) IUserRepo {
	return &UserRepo{
		app: deps,
	}
}

// FindAll retrieves users with filtering, sorting, and pagination
func (r *UserRepo) FindAll(ctx context.Context, params utils.QueryParams) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.app.GetDBWithContext(ctx).Model(&models.User{})
	query = query.Scopes(utils.ApplySearch(params.Search, "username", "name", "email", "phone"))
	query = query.Scopes(utils.ApplyFilter(params.Filters, allowedUserFilter))
	query = query.Scopes(utils.ApplyDateFilter(params.DateFrom, params.DateTo, params.DateField))

	// fetch total count for pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// apply sorting if enabled
	query = query.Scopes(utils.ApplySorting(params))

	// apply pagination if enabled
	query = query.Scopes(utils.ApplyPagination(params))

	// execute query to fetch users
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindById retrieves a user by ID
func (r *UserRepo) FindById(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := r.app.GetDBWithContext(ctx).First(&user, "id = ?", id).Error
	return &user, err
}

// Create adds a new user
func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	return r.app.GetDBWithContext(ctx).Create(user).Error
}

// Update user
func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	return r.app.GetDBWithContext(ctx).Save(user).Error
}

// UpdateById updates specific fields of a user by its ID
func (r *UserRepo) UpdateById(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.app.GetDBWithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete removes a user by ID
func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	return r.app.GetDBWithContext(ctx).Where("id = ?", id).Delete(&models.User{}).Error
}

// MarkAsActive marks a user as read
func (r *UserRepo) MarkAsActive(ctx context.Context, id int64) (bool, error) {
	err := r.app.GetDBWithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", true).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

// MarkAsInActive marks a user as read
func (r *UserRepo) MarkAsInActive(ctx context.Context, id int64) (bool, error) {
	err := r.app.GetDBWithContext(ctx).Model(&models.User{}).
		Where("id = ?", id).
		Update("is_active", false).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
