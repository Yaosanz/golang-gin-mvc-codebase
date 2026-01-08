package repositories

import (
	"context"

	"go-starter-app/app/http/utils"
	"go-starter-app/app/models"
)

// allowedShortenlinkFilter defines the fields that can be used for filtering shortenlink
var allowedShortenlinkFilter = map[string]bool{
	"created_by": true,
}

type ShortenlinkRepo struct {
	app IRepoDependencies
}

type IShortenlinkRepo interface {
	FindAll(ctx context.Context, params utils.QueryParams) ([]models.ShortenLink, int64, error)
	FindById(ctx context.Context, id string) (*models.ShortenLink, error)
	FindByCode(ctx context.Context, code string) (*models.ShortenLink, error)
	Create(ctx context.Context, data *models.ShortenLink) error
	Update(ctx context.Context, data *models.ShortenLink) error
	Delete(ctx context.Context, id string) error
}

func NewShortenlinkRepo(deps IRepoDependencies) IShortenlinkRepo {
	return &ShortenlinkRepo{
		app: deps,
	}
}

// ==========================
// FIND ALL (SAFE VERSION)
// ==========================
func (r *ShortenlinkRepo) FindAll(
	ctx context.Context,
	params utils.QueryParams,
) ([]models.ShortenLink, int64, error) {

	var data []models.ShortenLink
	var total int64

	db := r.app.GetDBWithContext(ctx).
		Model(&models.ShortenLink{})

	// filtering
	db = db.Scopes(
		utils.ApplySearch(params.Search, "original_url", "short_code"),
		utils.ApplyFilter(params.Filters, allowedShortenlinkFilter),
		utils.ApplyDateFilter(params.DateFrom, params.DateTo, params.DateField),
	)

	// count first
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ===== SAFE SORTING =====
	if params.SortBy != "" {
		db = db.Scopes(utils.ApplySorting(params))
	} else {
		// DEFAULT ORDER (WAJIB ADA)
		db = db.Order("created_at DESC")
	}

	// pagination
	db = db.Scopes(utils.ApplyPagination(params))

	// execute
	if err := db.Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, total, nil
}

// ==========================
// FIND BY ID
// ==========================
func (r *ShortenlinkRepo) FindById(ctx context.Context, id string) (*models.ShortenLink, error) {
	var data models.ShortenLink
	err := r.app.GetDBWithContext(ctx).
		First(&data, "id = ?", id).
		Error
	return &data, err
}

// ==========================
// FIND BY CODE
// ==========================
func (r *ShortenlinkRepo) FindByCode(ctx context.Context, code string) (*models.ShortenLink, error) {
	var data models.ShortenLink
	err := r.app.GetDBWithContext(ctx).
		Where("short_code = ?", code).
		First(&data).
		Error
	return &data, err
}

// ==========================
// CREATE
// ==========================
func (r *ShortenlinkRepo) Create(ctx context.Context, data *models.ShortenLink) error {
	return r.app.GetDBWithContext(ctx).Create(data).Error
}

// ==========================
// UPDATE
// ==========================
func (r *ShortenlinkRepo) Update(ctx context.Context, data *models.ShortenLink) error {
	return r.app.GetDBWithContext(ctx).Save(data).Error
}

// ==========================
// DELETE
// ==========================
func (r *ShortenlinkRepo) Delete(ctx context.Context, id string) error {
	return r.app.GetDBWithContext(ctx).
		Where("id = ?", id).
		Delete(&models.ShortenLink{}).
		Error
}
