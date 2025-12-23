package interfaces

import (
	"context"
	"go-starter-app/app/repositories"
	"go-starter-app/app/services"
	"go-starter-app/app/validation"
	"go-starter-app/config"
	"go-starter-app/pkg/google"
	"go-starter-app/pkg/scheduler"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// IAppDependencies app dependencies inversion
type IAppDependencies interface {
	GetConfig() *config.Config
	GetDB() *gorm.DB
	GetDBTransaction() *gorm.DB
	GetDBWithContext(ctx context.Context) *gorm.DB
	GetRepo() *repositories.RepoContainer
	GetService() *services.ServiceContainer
	GetValidator() *validation.AppValidator
	GetScheduler() *scheduler.Scheduler
	GetMinio() *minio.Client
	GetFcm() *google.FCM
}
