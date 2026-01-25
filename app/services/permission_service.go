package services

import (
	"context"

	"go-starter-app/app/repositories"
)

type PermissionService struct {
	permissionRepo repositories.IPermissionRepo
}

func NewPermissionService(deps IServiceDependencies) *PermissionService {
	return &PermissionService{
		permissionRepo: deps.GetRepo().PermissionRepo,
	}
}

func (s *PermissionService) HasPermission(ctx context.Context, roleName string, code string) (bool, error) {
	return s.permissionRepo.HasPermission(ctx, roleName, code)
}
