package services

import (
	"context"
	"errors"
	"time"

	"go-starter-app/app/http/dto"
	"go-starter-app/app/http/utils"
	"go-starter-app/app/models"
	"go-starter-app/helpers"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserService struct {
	app   IServiceDependencies
	cache helpers.CacheInterface
}

type IUserService interface {
	FindAll(ctx context.Context, params utils.QueryParams) ([]models.User, int64, error)
	FindById(ctx context.Context, id string) (*models.User, error)
	Create(ctx context.Context, dto *dto.CreateUserDTO) error
	Update(ctx context.Context, id string, dto *dto.UpdateUserDTO) error
	Delete(ctx context.Context, id string) error
}

func NewUserService(deps IServiceDependencies) *UserService {
	// Initialize cache (Redis if available, otherwise no-op)
	var cache helpers.CacheInterface
	if redisClient := deps.GetRedis(); redisClient != nil {
		redisCli := redisClient.(*redis.Client)
		cache = helpers.NewRedisClient(redisCli.Options().Addr, redisCli.Options().Password, redisCli.Options().DB)
	} else {
		cache = helpers.NewNoOpCache()
	}

	return &UserService{
		app:   deps,
		cache: cache,
	}
}

// FindAll retrieves users with filtering, sorting, and pagination
func (s *UserService) FindAll(ctx context.Context, params utils.QueryParams) ([]models.User, int64, error) {
	return s.app.GetRepo().UserRepo.FindAll(ctx, params)
}

// FindById retrieves a user by ID with optimized caching
func (s *UserService) FindById(ctx context.Context, id string) (*models.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid user id")
	}

	cacheKey := helpers.UserCacheKey(id)

	// Try cache first with lightweight cache data
	var cacheData models.UserCacheData
	err := s.cache.Get(ctx, cacheKey, &cacheData)
	if err == nil && cacheData.ID != "" {
		// Convert cache data back to User model
		userID, name, username, email, phone, isActive, role, roleID, err := cacheData.ToUserData()
		if err != nil {
			// Cache data corrupted, fall through to database
		} else {
			return &models.User{
				ID:       userID,
				Name:     name,
				Username: username,
				Email:    email,
				Phone:    phone,
				IsActive: isActive,
				Role:     role,
				RoleID:   roleID,
				// Note: Relationships not cached for performance
			}, nil
		}
	}

	// Cache miss or corrupted data, get from database
	userPtr, err := s.app.GetRepo().UserRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Cache the lightweight version for 5 minutes (reduced from 30 minutes)
	cacheData = *models.NewUserCacheData(
		userPtr.ID,
		userPtr.Name,
		userPtr.Username,
		userPtr.Email,
		userPtr.Phone,
		userPtr.IsActive,
		userPtr.GetPrimaryRole(),
		userPtr.GetPrimaryRoleID(),
	)
	s.cache.Set(ctx, cacheKey, cacheData, 5*time.Minute)

	return userPtr, nil
}

// Create adds a new user
func (s *UserService) Create(ctx context.Context, dto *dto.CreateUserDTO) error {
	// Validate unique constraints before creating
	var existingUser models.User

	// Check if username already exists
	if err := s.app.GetDB().Where("username = ?", dto.Username).First(&existingUser).Error; err == nil {
		return errors.New("username already exists")
	}

	// Check if email already exists
	if err := s.app.GetDB().Where("email = ?", dto.Email).First(&existingUser).Error; err == nil {
		return errors.New("email already exists")
	}

	// Check if phone already exists (if provided)
	if dto.Phone != nil && *dto.Phone != "" {
		if err := s.app.GetDB().Where("phone = ?", *dto.Phone).First(&existingUser).Error; err == nil {
			return errors.New("phone number already exists")
		}
	}

	// Get the default "user" role
	var userRole models.Role
	if err := s.app.GetDB().Where("name = ?", "user").First(&userRole).Error; err != nil {
		return errors.New("default user role not found")
	}

	user := &models.User{
		Name:     dto.Name,
		Username: dto.Username,
		Email:    dto.Email,
		Phone:    dto.Phone,
		Role:     models.RoleUser,
		RoleID:   userRole.ID,
		Roles:    []models.Role{userRole}, // Set the relationship
		IsActive: true, // New users are active by default
	}

	// Hash the password before saving
	password, err := helpers.HashPassword(dto.Password, 12)
	if err != nil {
		return err
	}
	user.Password = password

	// Create with transaction
	return helpers.RunInTransaction(ctx, s.app.GetDB(), func(ctx context.Context, tx *gorm.DB) error {
		// Create the user (GORM will automatically create the user_roles association due to the many2many relationship)
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		return nil
	})
}

// Update modifies an existing user
func (s *UserService) Update(ctx context.Context, id string, dto *dto.UpdateUserDTO) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid user id")
	}

	// Validate unique constraints before updating
	var existingUser models.User

	// Check if username already exists (excluding current user)
	if dto.Username != nil {
		if err := s.app.GetDB().Where("username = ? AND id != ?", *dto.Username, id).First(&existingUser).Error; err == nil {
			return errors.New("username already exists")
		}
	}

	// Check if email already exists (excluding current user)
	if dto.Email != nil {
		if err := s.app.GetDB().Where("email = ? AND id != ?", *dto.Email, id).First(&existingUser).Error; err == nil {
			return errors.New("email already exists")
		}
	}

	// Check if phone already exists (excluding current user)
	if dto.Phone != nil && *dto.Phone != "" {
		if err := s.app.GetDB().Where("phone = ? AND id != ?", *dto.Phone, id).First(&existingUser).Error; err == nil {
			return errors.New("phone number already exists")
		}
	}

	// Update with transaction
	err := helpers.RunInTransaction(ctx, s.app.GetDB(), func(ctx context.Context, tx *gorm.DB) error {
		user, err := s.app.GetRepo().UserRepo.FindById(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user not found")
			}
			return err
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
	})
	if err != nil {
		return err
	}

	// Invalidate cache
	s.cache.Delete(ctx, helpers.UserCacheKey(id))

	return nil
}

// Delete removes a user by ID
func (s *UserService) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid user id")
	}

	// Delete with transaction
	err := helpers.RunInTransaction(ctx, s.app.GetDB(), func(ctx context.Context, tx *gorm.DB) error {
		return s.app.GetRepo().UserRepo.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	// Invalidate cache
	s.cache.Delete(ctx, helpers.UserCacheKey(id))

	return nil
}
