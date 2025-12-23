package interfaces

import (
	"go-starter-app/app/repositories"
	"go-starter-app/app/services"
	"go-starter-app/app/validation"
	"go-starter-app/config"
	"go-starter-app/pkg/oca"
)

// KernelDependencies defines the dependencies required by the http kernel.
type KernelDependencies interface {
	GetConfig() *config.Config
	GetRepo() *repositories.RepoContainer
	GetService() *services.ServiceContainer
	GetValidator() *validation.AppValidator
	GetOca() *oca.Client
}
