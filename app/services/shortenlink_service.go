package services

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"go-starter-app/app/http/utils"
	"go-starter-app/app/models"
	"go-starter-app/app/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IShortenlinkService interface {
Create(ctx context.Context, originalURL, userID string) (*models.ShortenLink, error)
FindAll(ctx context.Context, userID string) ([]*models.ShortenLink, error)
FindByID(ctx context.Context, id, userID string) (*models.ShortenLink, error)
Update(ctx context.Context, id, originalURL, userID string) (*models.ShortenLink, error)
Delete(ctx context.Context, id, userID string) error
GetByCode(ctx context.Context, code string) (*models.ShortenLink, error)
}

type ShortenlinkService struct {
	repo repositories.IShortenlinkRepo
}

func NewShortenlinkService(
	repo repositories.IShortenlinkRepo,
) IShortenlinkService {
	return &ShortenlinkService{
		repo: repo,
	}
}

// ==========================
// CREATE SHORT LINK
// ==========================
func (s *ShortenlinkService) Create(
	ctx context.Context,
	originalURL string,
	userID string,
) (*models.ShortenLink, error) {

	// 🔒 Parse & validasi UUID dari JWT (user yang login)
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	// 🔑 Generate UUID baru untuk shortlink
	shortlinkID := uuid.New()

	now := time.Now()

	data := &models.ShortenLink{
		ID:          shortlinkID,       // UUID unik shortlink
		OriginalURL: originalURL,       // URL asli
		ShortCode:   generateShortCode(6), // Kode pendek 6 karakter
		UserID:      uid,               // Pemilik shortlink (dari JWT)
		CreatedBy:   uid,               // User yang membuat
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Simpan ke repository
	if err := s.repo.Create(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}


// ==========================
// FIND ALL SHORT LINKS
// ==========================
func (s *ShortenlinkService) FindAll(
	ctx context.Context,
	userID string,
) ([]*models.ShortenLink, error) {

	if _, err := uuid.Parse(userID); err != nil {
		return nil, errors.New("invalid user id")
	}

	params := &utils.QueryParams{
		Filters: map[string]string{
			"user_id": userID,
		},
	}
	data, _, err := s.repo.FindAll(ctx, *params)
	if err != nil {
		return nil, err
	}

	// Convert []models.ShortenLink to []*models.ShortenLink
	result := make([]*models.ShortenLink, len(data))
	for i := range data {
		result[i] = &data[i]
	}

	return result, nil
}

// ==========================
// FIND BY ID
// ==========================
func (s *ShortenlinkService) FindByID(
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

	data, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shortlink not found")
		}
		return nil, err
	}

	data.OriginalURL = originalURL
	data.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, data); err != nil {
		return nil, err
	}

	return data, nil
}

// ==========================
// GET BY SHORT CODE
// ==========================
func (s *ShortenlinkService) GetByCode(
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

	return s.repo.Delete(ctx, id)
}

// ==========================
// SHORT CODE GENERATOR
// ==========================
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateShortCode(length int) string {
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}
