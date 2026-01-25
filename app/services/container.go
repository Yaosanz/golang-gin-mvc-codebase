package services

import (
	"context"

	"go-starter-app/app/repositories"
	"go-starter-app/app/validation"
	"go-starter-app/config"
	"go-starter-app/pkg/google"
	"go-starter-app/pkg/oca"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// IServiceDependencies app dependencies inversion
type IServiceDependencies interface {
	GetConfig() *config.Config
	GetDB() *gorm.DB
	GetDBWithContext(ctx context.Context) *gorm.DB
	GetValidator() *validation.AppValidator
	GetRepo() *repositories.RepoContainer
	GetService() *ServiceContainer
	GetFcm() *google.FCM
	GetMinio() *minio.Client
	GetOca() *oca.Client
	GetRedis() *redis.Client
}

// SERVICE CONTAINER
type ServiceContainer struct {
	UserService        *UserService
	AuthService        *AuthService
	ShortenlinkService IShortenlinkService
	PermissionService  *PermissionService
}

func NewServiceContainer(deps IServiceDependencies) *ServiceContainer {
	cfg := deps.GetConfig().Jwt()

	return &ServiceContainer{
		UserService: NewUserService(deps),

		AuthService: NewAuthService(
			deps,
			deps.GetRepo().UserRepo,
			cfg.Secret,
			cfg.Issuer,
			cfg.ExpiredIn,
		),

		ShortenlinkService: NewShortenlinkService(deps),

		PermissionService: NewPermissionService(deps),
	}
}
