package services

import (
"context"
"errors"
"fmt"

"go-starter-app/app/http/dto"
"go-starter-app/app/http/utils"
"go-starter-app/app/models"
"go-starter-app/helpers"

"github.com/google/uuid"
"gorm.io/gorm"
)

type UserService struct {
app          IServiceDependencies
cache        helpers.CacheInterface
cacheManager *helpers.CacheManager
}

type IUserService interface {
FindAll(ctx context.Context, params utils.QueryParams) ([]models.User, int64, error)
FindById(ctx context.Context, id string) (*models.User, error)
FindByUsername(ctx context.Context, username string) (*models.User, error)
FindByEmail(ctx context.Context, email string) (*models.User, error)
Create(ctx context.Context, dto *dto.CreateUserDTO) error
Update(ctx context.Context, id string, dto *dto.UpdateUserDTO) error
Delete(ctx context.Context, id string) error
InvalidateUserCache(ctx context.Context, id string) error
}

func NewUserService(deps IServiceDependencies) *UserService {
var cache helpers.CacheInterface
if redisClient := deps.GetRedis(); redisClient != nil {
redisCache := helpers.NewRedisClient(
redisClient.Options().Addr,
redisClient.Options().Password,
redisClient.Options().DB,
)
cache = redisCache
} else {
cache = helpers.NewNoOpCache()
}

return &UserService{
app:          deps,
cache:        cache,
cacheManager: helpers.NewCacheManager(cache),
}
}

// FindAll retrieves users with filtering, sorting, and pagination
func (s *UserService) FindAll(ctx context.Context, params utils.QueryParams) ([]models.User, int64, error) {
return s.app.GetRepo().UserRepo.FindAll(ctx, params)
}

// FindById retrieves a user by ID with optimized caching strategy
func (s *UserService) FindById(ctx context.Context, id string) (*models.User, error) {
if _, err := uuid.Parse(id); err != nil {
return nil, fmt.Errorf("invalid user id: %w", err)
}

cacheKey := helpers.UserCacheKey(id)
var cacheData models.UserCacheData
err := s.cache.Get(ctx, cacheKey, &cacheData)
if err == nil && cacheData.ID != "" {
userID, name, username, email, phone, isActive, role, roleID, err := cacheData.ToUserData()
if err == nil {
return &models.User{
ID:       userID,
Name:     name,
Username: username,
Email:    email,
Phone:    phone,
IsActive: isActive,
Role:     role,
RoleID:   roleID,
}, nil
}
}

userPtr, err := s.app.GetRepo().UserRepo.FindById(ctx, id)
if err != nil {
if errors.Is(err, gorm.ErrRecordNotFound) {
return nil, errors.New("user not found")
}
return nil, fmt.Errorf("error finding user: %w", err)
}

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
_ = s.cache.Set(ctx, cacheKey, cacheData, helpers.UserCacheTTL)

return userPtr, nil
}

// FindByUsername retrieves a user by username with caching
func (s *UserService) FindByUsername(ctx context.Context, username string) (*models.User, error) {
if username == "" {
return nil, errors.New("username cannot be empty")
}

cacheKey := helpers.UserByUsernameCacheKey(username)
var user *models.User
err := s.cacheManager.GetOrSet(ctx, cacheKey, helpers.UserCacheTTL, func() (interface{}, error) {
return s.app.GetRepo().UserRepo.FindByUsername(ctx, username)
}, &user)

return user, err
}

// FindByEmail retrieves a user by email with caching
func (s *UserService) FindByEmail(ctx context.Context, email string) (*models.User, error) {
if email == "" {
return nil, errors.New("email cannot be empty")
}

cacheKey := helpers.UserByEmailCacheKey(email)
var user *models.User
err := s.cacheManager.GetOrSet(ctx, cacheKey, helpers.UserCacheTTL, func() (interface{}, error) {
return s.app.GetRepo().UserRepo.FindByEmail(ctx, email)
}, &user)

return user, err
}

// Create adds a new user with transaction
func (s *UserService) Create(ctx context.Context, dto *dto.CreateUserDTO) error {
var existingUser models.User

if err := s.app.GetDB().Where("username = ?", dto.Username).First(&existingUser).Error; err == nil {
return errors.New("username already exists")
}

if err := s.app.GetDB().Where("email = ?", dto.Email).First(&existingUser).Error; err == nil {
return errors.New("email already exists")
}

if dto.Phone != nil && *dto.Phone != "" {
if err := s.app.GetDB().Where("phone = ?", *dto.Phone).First(&existingUser).Error; err == nil {
return errors.New("phone number already exists")
}
}

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
Roles:    []models.Role{userRole},
IsActive: true,
}

password, err := helpers.HashPassword(dto.Password, 12)
if err != nil {
return fmt.Errorf("error hashing password: %w", err)
}
user.Password = password

return helpers.RunInTransaction(ctx, s.app.GetDB(), func(ctx context.Context, tx *gorm.DB) error {
if err := tx.Create(user).Error; err != nil {
return fmt.Errorf("error creating user: %w", err)
}
return nil
})
}

// Update modifies an existing user with cache invalidation
func (s *UserService) Update(ctx context.Context, id string, dto *dto.UpdateUserDTO) error {
if _, err := uuid.Parse(id); err != nil {
return errors.New("invalid user id")
}

var existingUser models.User

if dto.Username != nil {
if err := s.app.GetDB().Where("username = ? AND id != ?", *dto.Username, id).First(&existingUser).Error; err == nil {
return errors.New("username already exists")
}
}

if dto.Email != nil {
if err := s.app.GetDB().Where("email = ? AND id != ?", *dto.Email, id).First(&existingUser).Error; err == nil {
return errors.New("email already exists")
}
}

if dto.Phone != nil && *dto.Phone != "" {
if err := s.app.GetDB().Where("phone = ? AND id != ?", *dto.Phone, id).First(&existingUser).Error; err == nil {
return errors.New("phone number already exists")
}
}

err := helpers.RunInTransaction(ctx, s.app.GetDB(), func(ctx context.Context, tx *gorm.DB) error {
user, err := s.app.GetRepo().UserRepo.FindById(ctx, id)
if err != nil {
if errors.Is(err, gorm.ErrRecordNotFound) {
return errors.New("user not found")
}
return fmt.Errorf("error finding user: %w", err)
}

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
return fmt.Errorf("error hashing password: %w", err)
}
user.Password = password
}

return s.app.GetRepo().UserRepo.Update(ctx, user)
})
if err != nil {
return err
}

return s.InvalidateUserCache(ctx, id)
}

// Delete removes a user by ID with cache invalidation
func (s *UserService) Delete(ctx context.Context, id string) error {
if _, err := uuid.Parse(id); err != nil {
return errors.New("invalid user id")
}

err := helpers.RunInTransaction(ctx, s.app.GetDB(), func(ctx context.Context, tx *gorm.DB) error {
return s.app.GetRepo().UserRepo.Delete(ctx, id)
})
if err != nil {
return err
}

return s.InvalidateUserCache(ctx, id)
}

// InvalidateUserCache removes all cached data for a user
func (s *UserService) InvalidateUserCache(ctx context.Context, id string) error {
keys := []string{
helpers.UserCacheKey(id),
helpers.UserPermissionsCacheKey(id),
helpers.UserRolesCacheKey(id),
}

return s.cacheManager.Invalidate(ctx, keys...)
}
