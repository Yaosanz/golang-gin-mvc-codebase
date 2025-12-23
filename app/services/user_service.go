package services

import (
	"context"
	"go-starter-app/app/http/dto"
	"go-starter-app/app/http/utils"
	"go-starter-app/app/models"
	"go-starter-app/helpers"
)

type UserService struct {
	app IServiceDependencies
}

type IUserService interface {
	FindAll(ctx context.Context, params utils.QueryParams) ([]models.User, int64, error)
	FindById(ctx context.Context, id int64) (*models.User, error)
	Create(ctx context.Context, dto *dto.CreateUserDTO) error
	Update(ctx context.Context, dto *dto.UpdateUserDTO) error
	Delete(ctx context.Context, id int64) error
}

func NewUserService(deps IServiceDependencies) *UserService {
	return &UserService{
		app: deps,
	}
}

// FindAll retrieves users with filtering, sorting, and pagination
func (s *UserService) FindAll(ctx context.Context, params utils.QueryParams) ([]models.User, int64, error) {
	return s.app.GetRepo().UserRepo.FindAll(ctx, params)
}

// FindById retrieves a user by ID
func (s *UserService) FindById(ctx context.Context, id int64) (*models.User, error) {
	return s.app.GetRepo().UserRepo.FindById(ctx, id)
}

// Create adds a new user
func (s *UserService) Create(ctx context.Context, dto *dto.CreateUserDTO) error {
	user := &models.User{
		Name:     dto.Name,
		Username: dto.Username,
		Email:    dto.Email,
		Phone:    dto.Phone,
	}

	// Hash the password before saving
	password, err := helpers.HashPassword(dto.Password, 12)
	if err != nil {
		return err
	}
	user.Password = password

	return s.app.GetRepo().UserRepo.Create(ctx, user)
}

// Update modifies an existing user
func (s *UserService) Update(ctx context.Context, id int64, dto *dto.UpdateUserDTO) error {
	user, err := s.app.GetRepo().UserRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return nil // or return an error indicating user not found
	}

	// Update fields if they are provided
	if dto.Name != nil {
		user.Name = *dto.Name
	}
	if dto.Username != nil {
		user.Username = *dto.Username
	}
	if dto.Email != nil {
		user.Email = *dto.Email
	}
	if dto.Phone != nil {
		user.Phone = dto.Phone
	}
	if dto.Password != nil {
		password, err := helpers.HashPassword(*dto.Password, 12)
		if err != nil {
			return err
		}
		user.Password = password
	}

	return s.app.GetRepo().UserRepo.Update(ctx, user)
}

// Delete removes a user by ID
func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.app.GetRepo().UserRepo.Delete(ctx, id)
}
