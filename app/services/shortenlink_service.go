package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go-starter-app/app/http/utils"
	"go-starter-app/app/models"
	"go-starter-app/app/repositories"
	"go-starter-app/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IShortenlinkService interface {
	Create(ctx context.Context, originalURL, customCode, userID string) (*models.ShortenLink, error)
	FindAll(ctx context.Context, userID string) ([]*models.ShortenLink, error)
	FindById(ctx context.Context, id, userID string) (*models.ShortenLink, error)
	FindByCode(ctx context.Context, code string) (*models.ShortenLink, error)
	Update(ctx context.Context, id, originalURL, userID string) (*models.ShortenLink, error)
	UpdateByCode(ctx context.Context, code, originalURL, userID string) (*models.ShortenLink, error)
	Delete(ctx context.Context, id, userID string) error
	GetByCode(ctx context.Context, code string) (*models.ShortenLink, error)
	Redirect(ctx context.Context, code string) (string, error)
	InvalidateShortenLinkCache(ctx context.Context, code string) error
}

type ShortenlinkService struct {
	repo         repositories.IShortenlinkRepo
	cache        helpers.CacheInterface
	cacheManager *helpers.CacheManager
	deps         IServiceDependencies
}

func NewShortenlinkService(
	deps IServiceDependencies,
	repo repositories.IShortenlinkRepo,
) IShortenlinkService {
	// Initialize cache (Redis if available, otherwise no-op)
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

	return &ShortenlinkService{
		repo:         repo,
		cache:        cache,
		cacheManager: helpers.NewCacheManager(cache),
		deps:         deps,
	}
}

// ==========================
// CREATE SHORT LINK
// ==========================
func (s *ShortenlinkService) Create(
	ctx context.Context,
	originalURL string,
	customCode string,
	userID string,
) (*models.ShortenLink, error) {

	// Parse & validasi UUID dari JWT (user yang login)
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	// Generate UUID baru untuk shortlink
	shortlinkID := uuid.New()

	now := time.Now()

	// Use provided custom code if any, otherwise generate
	code := customCode
	if code == "" {
		code = generateShortCode()
	} else {
		if existing, err := s.repo.FindByCode(ctx, code); err == nil && existing != nil && existing.ID != uuid.Nil {
			return nil, errors.New("short code already exists")
		}
	}

	data := &models.ShortenLink{
		ID:          shortlinkID, // UUID unik shortlink
		OriginalURL: originalURL, // URL asli
		ShortCode:   code,        // Kode pendek
		UserID:      uid,         // Pemilik shortlink (dari JWT)
		CreatedBy:   uid,         // User yang membuat
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Simpan ke repository dengan transaction
	createdLink, err := helpers.RunInTransactionWithResult(ctx, s.deps.GetDB(), func(ctx context.Context, tx *gorm.DB) (*models.ShortenLink, error) {
		if err := s.repo.Create(ctx, data); err != nil {
			return nil, fmt.Errorf("error creating shortlink: %w", err)
		}
		return data, nil
	})
	if err != nil {
		return nil, err
	}

	// Invalidate user's shortlinks list cache after creating new link
	_ = s.cacheManager.Invalidate(ctx, helpers.UserShortenLinksCacheKey(userID))

	return createdLink, nil
}


// ==========================
// FIND ALL SHORT LINKS
// ==========================
func (s *ShortenlinkService) FindAll(
	ctx context.Context,
	userID string,
) ([]*models.ShortenLink, error) {

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	// Check if user has permission to read all shortenlinks
	canReadAll := false
	if svcContainer := s.deps.GetService(); svcContainer != nil {
		if authSvc := svcContainer.GetAuthService(); authSvc != nil {
			// Check if user has admin role - admin can see all
			hasAdmin, err := authSvc.CheckRole(ctx, uid, "admin")
			if err == nil && hasAdmin {
				canReadAll = true
			} else {
				// Check if user has CMS role or specific permission
				hasCMS, err := authSvc.CheckRole(ctx, uid, "cms")
				if err == nil && hasCMS {
					canReadAll = true
				} else {
					// Check specific permission for reading all
					hasPerm, err := authSvc.CheckPermission(ctx, uid, "shortenlink:read")
					if err == nil && hasPerm {
						canReadAll = true
					}
				}
			}
		}
	}

	var cacheKey string
	var filters map[string]string

	if canReadAll {
		// Admin/CMS can see all links
		cacheKey = "all_shortenlinks"
		filters = map[string]string{}
	} else {
		// Regular users see only their own
		cacheKey = helpers.UserShortenLinksCacheKey(userID)
		filters = map[string]string{
			"user_id": userID,
		}
	}

	// Try cache first
	var data []models.ShortenLink
	err = s.cache.Get(ctx, cacheKey, &data)
	if err == nil && len(data) > 0 {
		// Convert []models.ShortenLink to []*models.ShortenLink
		result := make([]*models.ShortenLink, len(data))
		for i := range data {
			result[i] = &data[i]
		}
		return result, nil
	}

	// Cache miss, get from database
	params := &utils.QueryParams{
		Filters: filters,
	}
	links, _, err := s.repo.FindAll(ctx, *params)
	if err != nil {
		return nil, err
	}

	// Convert []models.ShortenLink to []*models.ShortenLink
	result := make([]*models.ShortenLink, len(links))
	for i := range links {
		result[i] = &links[i]
	}

	// Cache the result for 30 minutes (cache the slice of structs)
	cacheData := make([]models.ShortenLink, len(links))
	copy(cacheData, links)
	s.cache.Set(ctx, cacheKey, cacheData, 30*time.Minute)

	return result, nil
}

// ==========================
// FIND BY ID
// ==========================
func (s *ShortenlinkService) FindById(
	ctx context.Context,
	id string,
	userID string,
) (*models.ShortenLink, error) {

	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid shortenlink id")
	}

	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("invalid user id")
	}

	data, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shortlink not found")
		}
		return nil, err
	}

	return data, nil
}

// ==========================
// FIND BY CODE
// ==========================
func (s *ShortenlinkService) FindByCode(
	ctx context.Context,
	code string,
) (*models.ShortenLink, error) {

	data, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shortlink not found")
		}
		return nil, err
	}

	return data, nil
}

// ==========================
// UPDATE SHORT LINK
// ==========================
func (s *ShortenlinkService) Update(
	ctx context.Context,
	id string,
	originalURL string,
	userID string,
) (*models.ShortenLink, error) {

	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("invalid shortenlink id")
	}

	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("invalid user id")
	}

	// Update dengan transaction
	updatedLink, err := helpers.RunInTransactionWithResult(ctx, s.deps.GetDB(), func(ctx context.Context, tx *gorm.DB) (*models.ShortenLink, error) {
		data, err := s.repo.FindById(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("shortlink not found")
			}
			return nil, fmt.Errorf("error finding shortlink: %w", err)
		}

		data.OriginalURL = originalURL
		data.UpdatedAt = time.Now()

		if err := s.repo.Update(ctx, data); err != nil {
			return nil, fmt.Errorf("error updating shortlink: %w", err)
		}

		return data, nil
	})
	if err != nil {
		return nil, err
	}

	// Invalidate caches after successful update
	_ = s.cacheManager.Invalidate(ctx,
		helpers.ShortenLinkByCodeCacheKey(updatedLink.ShortCode),
		helpers.UserShortenLinksCacheKey(userID),
	)
	s.cache.Delete(ctx, "all_shortenlinks")

	return updatedLink, nil
}

// ==========================
// UPDATE BY CODE
// ==========================
func (s *ShortenlinkService) UpdateByCode(
	ctx context.Context,
	code string,
	originalURL string,
	userID string,
) (*models.ShortenLink, error) {

	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("invalid user id")
	}

	updatedLink, err := helpers.RunInTransactionWithResult(ctx, s.deps.GetDB(), func(ctx context.Context, tx *gorm.DB) (*models.ShortenLink, error) {
		data, err := s.repo.FindByCode(ctx, code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("shortlink not found")
			}
			return nil, fmt.Errorf("error finding shortlink: %w", err)
		}

		data.OriginalURL = originalURL
		data.UpdatedAt = time.Now()

		if err := s.repo.Update(ctx, data); err != nil {
			return nil, fmt.Errorf("error updating shortlink: %w", err)
		}

		return data, nil
	})
	if err != nil {
		return nil, err
	}

	_ = s.InvalidateShortenLinkCache(ctx, updatedLink.ShortCode)
	return updatedLink, nil
}

// ==========================
// GET BY SHORT CODE
// ==========================
func (s *ShortenlinkService) GetByCode(
	ctx context.Context,
	code string,
) (*models.ShortenLink, error) {

	cacheKey := helpers.ShortenLinkByCodeCacheKey(code)

	// Try cache first
	var data models.ShortenLink
	err := s.cache.Get(ctx, cacheKey, &data)
	if err == nil && data.ID != uuid.Nil {
		return &data, nil
	}

	// Cache miss, get from database
	dataPtr, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shortlink not found")
		}
		return nil, fmt.Errorf("error finding shortlink by code: %w", err)
	}

	// Cache the result with professional TTL for shortlink codes (longer TTL for read-heavy operation)
	_ = s.cache.Set(ctx, cacheKey, *dataPtr, helpers.ShortenLinkCodeCacheTTL)

	return dataPtr, nil
}

// ==========================
// REDIRECT BY SHORT CODE
// ==========================
func (s *ShortenlinkService) Redirect(
	ctx context.Context,
	code string,
) (string, error) {
	link, err := s.GetByCode(ctx, code)
	if err != nil {
		return "", err
	}
	return link.OriginalURL, nil
}

// ==========================
// DELETE SHORT LINK
// ==========================
func (s *ShortenlinkService) Delete(
	ctx context.Context,
	id string,
	userID string,
) error {

	if _, err := uuid.Parse(id); err != nil {
		return errors.New("invalid shortenlink id")
	}

	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("invalid user id")
	}

	// Delete dengan transaction
	var linkToDelete *models.ShortenLink
	err := helpers.RunInTransaction(ctx, s.deps.GetDB(), func(ctx context.Context, tx *gorm.DB) error {
		// Get the link first to get the short code for cache invalidation
		link, err := s.repo.FindById(ctx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("shortlink not found")
			}
			return fmt.Errorf("error finding shortlink: %w", err)
		}
		linkToDelete = link

		// Delete from database
		if err := s.repo.Delete(ctx, id); err != nil {
			return fmt.Errorf("error deleting shortlink: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Invalidate caches after successful delete
	_ = s.cacheManager.Invalidate(ctx,
		helpers.ShortenLinkByCodeCacheKey(linkToDelete.ShortCode),
		helpers.UserShortenLinksCacheKey(userID),
	)

	return nil
}

// ==========================
// SHORT CODE GENERATOR
// ==========================
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateShortCode() string {
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}

// ==========================
// CACHE INVALIDATION
// ==========================

// InvalidateShortenLinkCache removes all cached data for a shortlink
// Called after updates/deletes to ensure data consistency
func (s *ShortenlinkService) InvalidateShortenLinkCache(ctx context.Context, code string) error {
	// Find the shortlink to get user ID for invalidating user's list
	link, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		// Log but don't fail - cache invalidation failure shouldn't break API
		fmt.Printf("Error finding shortlink for cache invalidation: %v\n", err)
		return nil
	}

	keys := []string{
		helpers.ShortenLinkByCodeCacheKey(code),
		helpers.UserShortenLinksCacheKey(link.UserID.String()),
	}

	return s.cacheManager.Invalidate(ctx, keys...)
}

