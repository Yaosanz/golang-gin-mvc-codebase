package interfaces

import (
	"context"
	"go-starter-app/app/http/dto"

	"github.com/google/uuid"
)

// IAuthService defines the authentication service interface
type IAuthService interface {
	// Authentication methods
	SecureLogin(ctx context.Context, username, password string) (string, string, interface{}, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, interface{}, error)
	Register(ctx context.Context, dto *dto.RegisterDTO) error

	// Authorization methods
	CheckPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
	CheckRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error)
	IsUserActive(ctx context.Context, userID uuid.UUID) (bool, error)

	// Middleware helpers
	RequirePermission(permission string) func(ctx context.Context, userID uuid.UUID) error
	RequireRole(role string) func(ctx context.Context, userID uuid.UUID) error
	RequireActiveUser() func(ctx context.Context, userID uuid.UUID) error

	// Session management
	InvalidateUserSession(ctx context.Context, sessionID string) error
	IsSessionValid(ctx context.Context, sessionID string) (bool, error)
	InvalidateUserAuthCache(ctx context.Context, username string) error
}

// LoginContext contains user context for server-side authorization
type LoginContext struct {
	UserID      uuid.UUID
	Username    string
	Roles       []string
	Permissions []string
	IsActive    bool
	SessionID   string
}
