package bootstrap

import (
	"context"
	"fmt"
	"go-starter-app/app/repositories"
	"go-starter-app/app/services"
	"go-starter-app/app/validation"
	"go-starter-app/config"
	"go-starter-app/interfaces"
	"go-starter-app/pkg/database"
	"go-starter-app/pkg/google"
	"go-starter-app/pkg/oca"
	"go-starter-app/pkg/scheduler"
	"go-starter-app/pkg/storage"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	config    *config.Config
	db        *gorm.DB
	redis     *redis.Client
	repos     *repositories.RepoContainer
	srvc      *services.ServiceContainer
	validator *validation.AppValidator
	scheduler *scheduler.Scheduler
	minio     *minio.Client
	fcm       *google.FCM
	oca       *oca.Client
}

// NewAppBootstrap initializes the application
func NewAppBootstrap() (*App, error) {
	a := &App{}

	// load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: \n%v", err)
	}
	a.config = cfg

	// initialize database connection
	db, dbErr := database.NewPostgres(cfg)
	if dbErr != nil {
		log.Fatalf("Failed to initialize database: \n%v", dbErr)
	}
	a.db = db

	// initialize redis connection if enabled
	if cfg.Redis().Enabled {
		a.redis = database.NewRedis(cfg)
	} else {
		log.Println("Redis disabled via config")
		a.redis = nil
	}

	// initialize minio client
	minioClient, err := storage.NewMinio(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio: %w", err)
	}
	a.minio = minioClient

	// initialize fcm client (optional)
	fcmCredentialPath := os.Getenv("FCM_CREDENTIAL_PATH")

	if fcmCredentialPath == "" {
		log.Println("[FCM] skipped: FCM_CREDENTIAL_PATH not set")
	} else {
		credPath := fcmCredentialPath
		fcm, err := google.NewFCM(credPath)

		if err != nil {
			log.Printf("[FCM] initialization failed: %v", err)
		} else {
			log.Println("[FCM] initialized successfully")
			a.fcm = fcm
		}
	}


	a.validator = validation.NewAppValidator(db) // initialize validator with db connection

	// initialize repositories and services after validator initialized
	a.repos = repositories.NewRepoContainer(a) // initialize repositories
	a.srvc = services.NewServiceContainer(a)   // initialize services
	a.scheduler = scheduler.New()              // initialize scheduler
	a.oca = oca.NewClient(cfg)                 // initialize oca client

	return a, nil
}

// GetConfig returns config instance
func (a *App) GetConfig() *config.Config {
	return a.config
}

// GetDB returns gorm db instance
func (a *App) GetDB() *gorm.DB {
	return a.db
}

// GetRedis returns redis client instance
func (a *App) GetRedis() interface{} {
	return a.redis
}

// GetDBTransaction returns gorm db instance with transaction
func (a *App) GetDBTransaction() *gorm.DB {
	return a.db.Begin()
}

// GetDBWithContext returns gorm db instance with context
func (a *App) GetDBWithContext(ctx context.Context) *gorm.DB {
	return a.db.WithContext(ctx)
}

// GetValidator returns app validation instance
func (a *App) GetValidator() *validation.AppValidator {
	return a.validator
}

// GetRepo returns app repositories container instance
func (a *App) GetRepo() *repositories.RepoContainer {
	return a.repos
}

// GetService returns app services container instance
func (a *App) GetService() interfaces.IServiceContainer {
	return a.srvc
}

// GetScheduler returns scheduler instance
func (a *App) GetScheduler() *scheduler.Scheduler {
	return a.scheduler
}

// GetFcm returns fcm instance
func (a *App) GetFcm() *google.FCM {
	return a.fcm
}

// GetMinio returns minio instance
func (a *App) GetMinio() *minio.Client {
	return a.minio
}

// GetOca returns oca instance
func (a *App) GetOca() *oca.Client {
	return a.oca
}
