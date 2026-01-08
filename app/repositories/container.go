package repositories

import (
	"context"

	"gorm.io/gorm"
)

// IRepoDependencies app dependencies inversion
type IRepoDependencies interface {
	GetDB() *gorm.DB
	GetDBWithContext(ctx context.Context) *gorm.DB
}

type RepoContainer struct {
	UserRepo IUserRepo
	ShortenlinkRepo IShortenlinkRepo
	PermissionRepo IPermissionRepo
	// Add other repositories here as needed (using interfaces)
}

func NewRepoContainer(deps IRepoDependencies) *RepoContainer {
	return &RepoContainer{
		UserRepo: NewUserRepo(deps),
		ShortenlinkRepo: NewShortenlinkRepo(deps),
		PermissionRepo: NewPermissionRepo(deps),
		// Add other repositories here as needed
	}
}
