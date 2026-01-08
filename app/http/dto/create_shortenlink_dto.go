package dto

import "time"

// ==========================
// CREATE SHORTENLINK DTO
// ==========================
type CreateShortenLinkDTO struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
}


type UpdateShortenLinkDTO struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
}

// ==========================
// RESPONSE SHORTENLINK DTO
// ==========================
type ResponseShortenLinkDTO struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ShortCode   string    `json:"short_code" example:"AbC123"`
	OriginalURL string    `json:"original_url" example:"https://www.google.com"`
	CreatedAt   time.Time `json:"created_at" example:"2026-01-01T10:00:00Z"`
}