package services

import (
	"context"
	"go-starter-app/app/repositories"
	"go-starter-app/app/validation"
	"go-starter-app/config"
	"go-starter-app/pkg/google"
	"go-starter-app/pkg/oca"

	"github.com/minio/minio-go/v7"
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
}

type ServiceContainer struct {
	UserService *UserService
	// Add other services here
}

func NewServiceContainer(deps IServiceDependencies) *ServiceContainer {
	return &ServiceContainer{
		UserService: NewUserService(deps),
		// Initialize other services here
	}
}
