package services

import (
	"context"
	"errors"
	"time"

	"go-starter-app/app/http/dto"
	"go-starter-app/app/models"
	"go-starter-app/app/repositories"
	"go-starter-app/helpers"

	"github.com/google/uuid"
)

// SecureAuthService provides ultra-fast, secure authentication following JWT best practices
type SecureAuthService struct {
	userRepo       repositories.IUserRepo
	roleRepo       repositories.IRoleRepo
	permissionRepo repositories.IPermissionRepo
	jwtSecret      string
	jwtIssuer      string
	jwtExpired     time.Duration
	refreshTTL     time.Duration
	deps           IServiceDependencies
	cache          helpers.CacheInterface
	authCache      *helpers.AuthCacheService
	authzService   *helpers.AuthorizationService
}

// NewSecureAuthService creates a secure auth service
func NewSecureAuthService(
	deps IServiceDependencies,
	userRepo repositories.IUserRepo,
	roleRepo repositories.IRoleRepo,
	permissionRepo repositories.IPermissionRepo,
	jwtSecret string,
	jwtIssuer string,
	jwtExpired time.Duration,
	cache helpers.CacheInterface,
) *SecureAuthService {
	return &SecureAuthService{
		deps:           deps,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		jwtSecret:      jwtSecret,
		jwtIssuer:      jwtIssuer,
		jwtExpired:     jwtExpired,
		cache:          cache,
		refreshTTL:     30 * 24 * time.Hour, // 30 days for refresh token
		authCache:      helpers.NewAuthCacheService(cache),
		authzService:   helpers.NewAuthorizationService(cache),
	}
}

// Register creates a new user account
func (s *SecureAuthService) Register(ctx context.Context, req *dto.RegisterDTO) error {
	createUserDTO := &dto.CreateUserDTO{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}
	return s.deps.GetService().GetUserService().Create(ctx, createUserDTO)
}

// SecureLogin performs ultra-fast, secure authentication
// Returns access token + refresh token + user context for authorization
func (s *SecureAuthService) SecureLogin(ctx context.Context, username, password string) (string, string, interface{}, error) {
	// Try cache first - ultra fast path
	authData, err := s.authCache.GetUserAuth(ctx, username)
	if err == nil {
		// Cache hit - verify password and generate token
		if !helpers.CheckPassword(password, authData.PasswordHash) {
			return "", "", nil, errors.New("invalid credentials")
		}

		if !authData.IsActive {
			return "", "", nil, errors.New("account is inactive")
		}

		// Generate secure JWT tokens and context
		accessToken, refreshToken, loginCtx, err := s.issueTokens(ctx, authData)
		if err != nil {
			return "", "", nil, err
		}

		return accessToken, refreshToken, loginCtx, nil
	}

	// Cache miss - fallback to database (one-time cost)
	return s.databaseSecureLogin(ctx, username, password)
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

// databaseSecureLogin handles login when cache miss occurs
func (s *SecureAuthService) databaseSecureLogin(ctx context.Context, username, password string) (string, string, interface{}, error) {
	// Get user from database
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", "", nil, errors.New("invalid credentials")
	}

	// Verify password
	if !helpers.CheckPassword(password, user.Password) {
		return "", "", nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", "", nil, errors.New("account is inactive")
	}

	// Get roles and permissions for caching and context
	userWithRoles, err := s.getUserWithRoles(ctx, user.ID)
	if err != nil {
		return "", "", nil, err
	}

	roles := make([]string, len(userWithRoles.Roles))
	for i, role := range userWithRoles.Roles {
		roles[i] = role.Name
	}

	permissions, err := s.getUserPermissions(ctx, user.ID)
	if err != nil {
		return "", "", nil, err
	}

	// Cache auth data for future fast logins
	s.authCache.CacheUserAuth(ctx, user.ID, user.Username, user.Email, user.Password, user.IsActive, roles, permissions)

	// Also cache permissions separately for authorization service
	s.authzService.CacheUserPermissions(ctx, user.ID, permissions)
	s.authzService.CacheUserRoles(ctx, user.ID, roles)

	// Cache user data for IsUserActive check
	userCacheData := &models.UserCacheData{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Phone:    "",
		IsActive: user.IsActive,
		Role:     user.Role,
		RoleID:   user.RoleID.String(),
	}
	cacheKey := helpers.UserCacheKey(user.ID.String())
	s.cache.Set(ctx, cacheKey, userCacheData, 10*time.Minute)

	// Create auth data for token generation
	authData := &helpers.UserAuthData{
		UserID:       user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.Password,
		IsActive:     user.IsActive,
		Roles:        roles,
		Permissions:  permissions,
	}

	// Generate secure tokens and context
	accessToken, refreshToken, loginCtx, err := s.issueTokens(ctx, authData)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, loginCtx, nil
}

// Refresh exchanges a refresh token for new access & refresh tokens (rotation)
func (s *SecureAuthService) Refresh(ctx context.Context, refreshToken string) (string, string, interface{}, error) {
	if refreshToken == "" {
		return "", "", nil, errors.New("refresh token required")
	}

	// Get session ID from cache (validates refresh token exists)
	oldSessionID, err := s.getRefreshTokenData(ctx, refreshToken)
	if err != nil {
		return "", "", nil, errors.New("invalid or expired refresh token")
	}

	// Get user ID from cache before deleting
	cacheKey := s.refreshTokenKey(refreshToken) + ":uid"
	var userIDStr string
	if err := s.cache.Get(ctx, cacheKey, &userIDStr); err != nil {
		return "", "", nil, errors.New("user context not found")
	}

	// One-time use: remove old token and blacklist previous session
	_ = s.cache.Delete(ctx, s.refreshTokenKey(refreshToken))
	_ = s.cache.Delete(ctx, cacheKey) // Also delete user ID cache
	if oldSessionID != "" {
		_ = s.InvalidateUserSession(ctx, oldSessionID)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", nil, errors.New("invalid user identifier")
	}

	user, err := s.userRepo.FindById(ctx, userID.String())
	if err != nil {
		return "", "", nil, errors.New("user not found")
	}
	if !user.IsActive {
		return "", "", nil, errors.New("account is inactive")
	}

	userWithRoles, err := s.getUserWithRoles(ctx, user.ID)
	if err != nil {
		return "", "", nil, err
	}

	roles := make([]string, len(userWithRoles.Roles))
	for i, role := range userWithRoles.Roles {
		roles[i] = role.Name
	}

	permissions, err := s.getUserPermissions(ctx, user.ID)
	if err != nil {
		return "", "", nil, err
	}

	// Refresh caches so authorization remains fast
	_ = s.authCache.CacheUserAuth(ctx, user.ID, user.Username, user.Email, user.Password, user.IsActive, roles, permissions)
	_ = s.authzService.CacheUserPermissions(ctx, user.ID, permissions)
	_ = s.authzService.CacheUserRoles(ctx, user.ID, roles)

	authData := &helpers.UserAuthData{
		UserID:       user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.Password,
		IsActive:     user.IsActive,
		Roles:        roles,
		Permissions:  permissions,
	}

	accessToken, newRefreshToken, loginCtx, err := s.issueTokens(ctx, authData)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, newRefreshToken, loginCtx, nil
}

// issueTokens creates both access and refresh tokens with session context
func (s *SecureAuthService) issueTokens(ctx context.Context, authData *helpers.UserAuthData) (string, string, *LoginContext, error) {
	userID, _, _, _, _, _, _, err := authData.ToAuthData()
	if err != nil {
		return "", "", nil, err
	}

	// Generate unique session ID for this login
	sessionID := helpers.GenerateSessionID()

	// Create minimal, secure JWT token - NO sensitive data exposed
	accessToken, err := helpers.GenerateSecureUserToken(
		userID,
		"user", // token type
		sessionID,
		authData.Roles,
		s.jwtSecret,
		s.jwtExpired,
		s.jwtIssuer,
	)
	if err != nil {
		return "", "", nil, err
	}

	// Generate and store refresh token
	refreshToken, err := s.generateAndStoreRefreshToken(ctx, &LoginContext{
		UserID:      userID,
		Username:    authData.Username,
		Roles:       authData.Roles,
		Permissions: authData.Permissions,
		IsActive:    authData.IsActive,
		SessionID:   sessionID,
	})
	if err != nil {
		return "", "", nil, err
	}

	// Return context for authorization
	loginCtx := &LoginContext{
		UserID:      userID,
		Username:    authData.Username,
		Roles:       authData.Roles,
		Permissions: authData.Permissions,
		IsActive:    authData.IsActive,
		SessionID:   sessionID,
	}

	return accessToken, refreshToken, loginCtx, nil
}

// generateAndStoreRefreshToken creates rotating refresh token and stores server-side
func (s *SecureAuthService) generateAndStoreRefreshToken(ctx context.Context, loginCtx *LoginContext) (string, error) {
	rt, err := helpers.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	// Store refresh token data: sessionID
	key := s.refreshTokenKey(rt)
	if err := s.cache.Set(ctx, key, loginCtx.SessionID, s.refreshTTL); err != nil {
		return "", err
	}

	// Also cache user ID for refresh endpoint
	if err := s.cache.Set(ctx, key+":uid", loginCtx.UserID.String(), s.refreshTTL); err != nil {
		return "", err
	}

	return rt, nil
}

// getRefreshTokenData retrieves session ID from cached refresh token
func (s *SecureAuthService) getRefreshTokenData(ctx context.Context, token string) (string, error) {
	var sessionID string
	if err := s.cache.Get(ctx, s.refreshTokenKey(token), &sessionID); err != nil {
		return "", err
	}
	return sessionID, nil
}

// refreshTokenKey generates cache key for refresh token
func (s *SecureAuthService) refreshTokenKey(token string) string {
	return "auth:refresh:" + token
}

// generateSecureToken creates a secure JWT token with minimal payload
func (s *SecureAuthService) generateSecureToken(authData *helpers.UserAuthData) (string, interface{}, error) {
	userID, _, _, _, _, _, _, err := authData.ToAuthData()
	if err != nil {
		return "", nil, err
	}

	// Generate unique session ID for this login
	sessionID := helpers.GenerateSessionID()

	// Create minimal, secure JWT token - NO sensitive data exposed
	token, err := helpers.GenerateSecureUserToken(
		userID,
		"user", // token type
		sessionID,
		authData.Roles,
		s.jwtSecret,
		s.jwtExpired,
		s.jwtIssuer,
	)
	if err != nil {
		return "", nil, err
	}

	// Return token + server-side context for authorization
	loginCtx := &LoginContext{
		UserID:      userID,
		Username:    authData.Username,
		Roles:       authData.Roles,
		Permissions: authData.Permissions,
		IsActive:    authData.IsActive,
		SessionID:   sessionID,
	}

	return token, loginCtx, nil
}

// Authorization Methods - Server-side only (secure)

// CheckPermission verifies if user has permission
func (s *SecureAuthService) CheckPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	return s.authzService.CheckPermission(ctx, userID, permission)
}

// CheckRole verifies if user has role
func (s *SecureAuthService) CheckRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error) {
	return s.authzService.CheckRole(ctx, userID, roleName)
}

// GetUserPermissions returns all user permissions
func (s *SecureAuthService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.authzService.GetUserPermissions(ctx, userID)
}

// GetUserRoles returns all user roles
func (s *SecureAuthService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	return s.authzService.GetUserRoles(ctx, userID)
}

// GetUserPermissionsFromDB loads permissions from database and repopulates cache
func (s *SecureAuthService) GetUserPermissionsFromDB(ctx context.Context, userID uuid.UUID) ([]string, error) {
	permissions, err := s.getUserPermissions(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Repopulate cache
	s.authzService.CacheUserPermissions(ctx, userID, permissions)

	// Also cache roles while we're at it
	user, err := s.getUserWithRoles(ctx, userID)
	if err == nil {
		roles := make([]string, len(user.Roles))
		for i, role := range user.Roles {
			roles[i] = role.Name
		}
		s.authzService.CacheUserRoles(ctx, userID, roles)
	}

	return permissions, nil
}

// IsUserActive checks if user account is active with database fallback
func (s *SecureAuthService) IsUserActive(ctx context.Context, userID uuid.UUID) (bool, error) {
	// Try cache first
	isActive, err := s.authzService.IsUserActive(ctx, userID)
	if err == nil {
		return isActive, nil
	}

	// Cache miss - fallback to database
	user, err := s.userRepo.FindById(ctx, userID.String())
	if err != nil {
		return false, err
	}

	// Repopulate cache for future checks
	userCacheData := &models.UserCacheData{
		ID:       user.ID.String(),
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		IsActive: user.IsActive,
	}
	cacheKey := helpers.UserCacheKey(user.ID.String())
	s.cache.Set(ctx, cacheKey, userCacheData, 30*time.Minute)

	return user.IsActive, nil
}

// Middleware Helpers for Secure Authorization

// RequirePermission creates a permission check middleware
func (s *SecureAuthService) RequirePermission(permission string) func(ctx context.Context, userID uuid.UUID) error {
	return helpers.RequirePermission(s.authzService, permission)
}

// RequireRole creates a role check middleware
func (s *SecureAuthService) RequireRole(role string) func(ctx context.Context, userID uuid.UUID) error {
	return helpers.RequireRole(s.authzService, role)
}

// RequireActiveUser creates an active user check middleware
func (s *SecureAuthService) RequireActiveUser() func(ctx context.Context, userID uuid.UUID) error {
	return helpers.RequireActiveUser(s.authzService)
}

// Session Management

// InvalidateUserSession invalidates a specific user session
func (s *SecureAuthService) InvalidateUserSession(ctx context.Context, sessionID string) error {
	// In Redis, we can store session blacklist or use JWT ID for revocation
	cacheKey := "session:blacklist:" + sessionID
	return s.cache.Set(ctx, cacheKey, true, 24*time.Hour) // Blacklist for 24 hours
}

// IsSessionValid checks if a session is still valid
func (s *SecureAuthService) IsSessionValid(ctx context.Context, sessionID string) (bool, error) {
	cacheKey := "session:blacklist:" + sessionID
	var blacklisted bool
	err := s.cache.Get(ctx, cacheKey, &blacklisted)
	if err != nil {
		// If not in cache, session is valid
		return true, nil
	}
	return !blacklisted, nil
}

// InvalidateUserAuthCache invalidates auth cache for a user
func (s *SecureAuthService) InvalidateUserAuthCache(ctx context.Context, username string) error {
	return s.authCache.InvalidateUserAuth(ctx, username)
}

// Helper methods (same as optimized version)

func (s *SecureAuthService) getUserWithRoles(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	userPtr, err := s.userRepo.FindById(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	if err := s.deps.GetDB().Preload("Roles").First(userPtr, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return userPtr, nil
}

func (s *SecureAuthService) getUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	user, err := s.getUserWithRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	permissions := make([]string, 0)
	permissionMap := make(map[string]bool)

	for _, role := range user.Roles {
		rolePermissions, err := s.permissionRepo.GetPermissionsByRoleID(ctx, role.ID)
		if err != nil {
			continue
		}

		for _, perm := range rolePermissions {
			if !permissionMap[perm.Name] {
				permissionMap[perm.Name] = true
				permissions = append(permissions, perm.Name)
			}
		}
	}

	return permissions, nil
}
