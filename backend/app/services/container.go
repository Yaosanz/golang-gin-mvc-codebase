package services

import (
	"context"
	"go-starter-app/interfaces"

	"go-starter-app/app/repositories"
	"go-starter-app/app/validation"
	"go-starter-app/config"
	"go-starter-app/helpers"
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
	GetService() interfaces.IServiceContainer
	GetFcm() *google.FCM
	GetMinio() *minio.Client
	GetOca() *oca.Client
	GetRedis() *redis.Client
}

// SERVICE CONTAINER
type ServiceContainer struct {
	UserService        *UserService
	AuthService        interfaces.IAuthService
	ShortenlinkService IShortenlinkService
	PermissionService  *PermissionService
}

func NewServiceContainer(deps IServiceDependencies) *ServiceContainer {
	cfg := deps.GetConfig().Jwt()

	// Initialize cache (Redis if available, otherwise no-op)
	var cache helpers.CacheInterface
	if redisClient := deps.GetRedis(); redisClient != nil {
		cache = helpers.NewRedisClient(redisClient.Options().Addr, redisClient.Options().Password, redisClient.Options().DB)
	} else {
		cache = helpers.NewNoOpCache()
	}

	return &ServiceContainer{
		UserService: NewUserService(deps),

		AuthService: NewSecureAuthService(
			deps,
			deps.GetRepo().UserRepo,
			deps.GetRepo().RoleRepo,
			deps.GetRepo().PermissionRepo,
			cfg.Secret,
			cfg.Issuer,
			cfg.ExpiredIn,
			cache,
		),

		ShortenlinkService: NewShortenlinkService(
			deps,
			deps.GetRepo().ShortenlinkRepo,
		),

		PermissionService: NewPermissionService(deps),
	}
}
func (sc *ServiceContainer) GetAuthService() interfaces.IAuthService {
	return sc.AuthService
}

func (sc *ServiceContainer) GetUserService() interfaces.IUserService {
	return sc.UserService
}

func (sc *ServiceContainer) GetShortenlinkService() interfaces.IShortenlinkService {
	return sc.ShortenlinkService
}

func (sc *ServiceContainer) GetPermissionService() interfaces.IPermissionService {
	return sc.PermissionService
}
