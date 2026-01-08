package services

import (
	"context"

	"go-starter-app/app/repositories"

	"github.com/google/uuid"
)

type PermissionService struct {
	permissionRepo repositories.IPermissionRepo
}

func NewPermissionService(deps IServiceDependencies) *PermissionService {
	return &PermissionService{
		permissionRepo: deps.GetRepo().PermissionRepo,
	}
}

func (s *PermissionService) HasPermission(ctx context.Context, roleID uuid.UUID, code string) (bool, error) {
	return s.permissionRepo.HasPermission(ctx, roleID, code)
}
