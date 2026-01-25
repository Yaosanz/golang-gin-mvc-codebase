package services

import (
	"context"
	"errors"
	"time"

	"go-starter-app/app/http/dto"
	"go-starter-app/app/models"
	"go-starter-app/app/repositories"
	"go-starter-app/helpers"
)

type AuthService struct {
	userRepo   repositories.IUserRepo
	jwtSecret  string
	jwtIssuer  string
	jwtExpired time.Duration
	deps       IServiceDependencies
}

func NewAuthService(
	deps IServiceDependencies,
	userRepo repositories.IUserRepo,
	jwtSecret string,
	jwtIssuer string,
	jwtExpired time.Duration,
) *AuthService {
	return &AuthService{
		deps:       deps,
		userRepo:   userRepo,
		jwtSecret: jwtSecret,
		jwtIssuer: jwtIssuer,
		jwtExpired: jwtExpired,
	}
}

// =======================
// LOGIN
// =======================
func (s *AuthService) Login(
	ctx context.Context,
	username, password string,
) (string, error) {

	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !helpers.CheckPassword(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	sessionID := helpers.GenerateSessionID()
	token, err := helpers.GenerateSecureUserToken(
		user.ID,
		"user",
		sessionID,
		[]string{user.Role},
		s.jwtSecret,
		s.jwtExpired,
		s.jwtIssuer,
	)
	if err != nil {
		return "", err
	}

	return token, nil
}

// =======================
// REGISTER
// =======================
func (s *AuthService) Register(
	ctx context.Context,
	req *dto.RegisterDTO,
) error {

	// cek username
	if _, err := s.userRepo.FindByUsername(ctx, req.Username); err == nil {
		return errors.New("username already exists")
	}

	// cek email
	if _, err := s.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := helpers.HashPassword(req.Password, 12)
	if err != nil {
		return err
	}

	// Get user role ID
	var userRole models.Role
	if err := s.deps.GetDB().Where("name = ?", "user").First(&userRole).Error; err != nil {
		return err
	}

	user := &models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: hashedPassword,
		Role:     models.RoleUser,
		RoleID:   userRole.ID,
		IsActive: true,
	}

	return s.userRepo.Create(ctx, user)
}
