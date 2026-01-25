package interfaces

import (
	"context"
	"go-starter-app/app/http/dto"
	"go-starter-app/app/http/utils"
	"go-starter-app/app/models"
	"go-starter-app/app/repositories"
	"go-starter-app/app/validation"
	"go-starter-app/config"
	"go-starter-app/pkg/google"
	"go-starter-app/pkg/scheduler"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type IUserService interface {
	FindAll(ctx context.Context, params utils.QueryParams) ([]models.User, int64, error)
	FindById(ctx context.Context, id string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, dto *dto.CreateUserDTO) error
	Update(ctx context.Context, id string, dto *dto.UpdateUserDTO) error
	Delete(ctx context.Context, id string) error
}

type IShortenlinkService interface {
	Create(ctx context.Context, originalURL, customCode, userID string) (*models.ShortenLink, error)
	FindAll(ctx context.Context, userID string) ([]*models.ShortenLink, error)
	FindById(ctx context.Context, id, userID string) (*models.ShortenLink, error)
	FindByCode(ctx context.Context, code string) (*models.ShortenLink, error)
	GetByCode(ctx context.Context, code string) (*models.ShortenLink, error)
	Update(ctx context.Context, id, originalURL, userID string) (*models.ShortenLink, error)
	UpdateByCode(ctx context.Context, code, originalURL, userID string) (*models.ShortenLink, error)
	Delete(ctx context.Context, id, userID string) error
	Redirect(ctx context.Context, code string) (string, error)
}

type IPermissionService interface {
	HasPermission(ctx context.Context, roleName string, code string) (bool, error)
}

// IServiceContainer defines the service container interface
type IServiceContainer interface {
	GetAuthService() IAuthService
	GetUserService() IUserService
	GetShortenlinkService() IShortenlinkService
	GetPermissionService() IPermissionService
	// Add other service getters as needed
}

// IAppDependencies app dependencies inversion
type IAppDependencies interface {
	GetConfig() *config.Config
	GetDB() *gorm.DB
	GetDBTransaction() *gorm.DB
	GetDBWithContext(ctx context.Context) *gorm.DB
	GetRepo() *repositories.RepoContainer
	GetService() IServiceContainer
	GetValidator() *validation.AppValidator
	GetScheduler() *scheduler.Scheduler
	GetMinio() *minio.Client
	GetFcm() *google.FCM
}
