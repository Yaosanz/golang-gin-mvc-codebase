package v1

import (
	"errors"
	"net/http"

	"go-starter-app/app/http/dto"
	"go-starter-app/app/http/utils"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShortenlinkController struct {
	app interfaces.KernelDependencies
}

func NewShortenlinkController(app interfaces.KernelDependencies) *ShortenlinkController {
	return &ShortenlinkController{
		app: app,
	}
}

//
// ==========================
// CREATE SHORTLINK (JWT)
// ==========================
func (c *ShortenlinkController) Create(ctx *gin.Context) {
	var req dto.CreateShortenLinkDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	if validationErrors := c.app.GetValidator().ValidateStruct(&req); validationErrors != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Validation error", validationErrors)
		return
	}

	userIDVal, ok := ctx.Get("user_id")
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userIDUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}
	userID := userIDUUID.String()

	data, err := c.app.GetService().
		GetShortenlinkService().
		Create(ctx.Request.Context(), req.OriginalURL, req.CustomCode, userID)

	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	status := http.StatusCreated
	utils.SendOne(ctx, gin.H{
		"id":   data.ID,
		"code": data.ShortCode,
		"url":  data.OriginalURL,
	}, http.StatusText(http.StatusCreated), &status)
}

//
// ==========================
// GET ALL SHORTLINKS (JWT)
// ==========================
func (c *ShortenlinkController) FindAll(ctx *gin.Context) {
	userIDVal, ok := ctx.Get("user_id")
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userIDUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}
	userID := userIDUUID.String()

	data, err := c.app.GetService().
		GetShortenlinkService().
		FindAll(ctx.Request.Context(), userID)
	if err != nil {
		ctx.Error(err) // Log the error
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SendOne(ctx, data, http.StatusText(http.StatusOK), nil)
}

//
// ==========================
// GET SHORTLINK BY ID (JWT)
// ==========================
func (c *ShortenlinkController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	userIDVal, ok := ctx.Get("user_id")
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userIDUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}
	userID := userIDUUID.String()

	data, err := c.app.GetService().
		GetShortenlinkService().
		FindById(ctx.Request.Context(), id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(ctx, http.StatusNotFound, "Shortlink not found", nil)
			return
		}
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SendOne(ctx, data, http.StatusText(http.StatusOK), nil)
}

//
// ==========================
// UPDATE SHORTLINK (JWT)
// ==========================
func (c *ShortenlinkController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	var req dto.UpdateShortenLinkDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	userIDVal, ok := ctx.Get("user_id")
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userIDUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}
	userID := userIDUUID.String()

	updated, err := c.app.GetService().
		GetShortenlinkService().
		Update(ctx.Request.Context(), id, req.OriginalURL, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(ctx, http.StatusNotFound, "Shortlink not found", nil)
			return
		}
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SendOne(ctx, gin.H{
		"id":   updated.ID,
		"code": updated.ShortCode,
		"url":  updated.OriginalURL,
	}, http.StatusText(http.StatusOK), nil)
}

//
// ==========================
// DELETE SHORTLINK (JWT)
// ==========================
func (c *ShortenlinkController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid ID", nil)
		return
	}

	userIDVal, ok := ctx.Get("user_id")
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userIDUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}
	userID := userIDUUID.String()

	err := c.app.GetService().
		GetShortenlinkService().
		Delete(ctx.Request.Context(), id, userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(ctx, http.StatusNotFound, "Shortlink not found", nil)
			return
		}
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SendOne(ctx, nil, http.StatusText(http.StatusOK), nil)
}

// ===== Compatibility endpoints used by Postman collection (code-based operations) =====

// POST /api/shortenlinks
func (c *ShortenlinkController) CreateByCode(ctx *gin.Context) {
	var req dto.CreateShortenLinkDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	userIDVal, ok := ctx.Get("user_id")
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userIDUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}

	created, err := c.app.GetService().GetShortenlinkService().Create(ctx.Request.Context(), req.OriginalURL, req.CustomCode, userIDUUID.String())
	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	status := http.StatusCreated
	utils.SendOne(ctx, gin.H{
		"code": created.ShortCode,
		"url":  created.OriginalURL,
		"id":   created.ID,
	}, "Shortlink created", &status)
}

// GET /api/shortenlinks/:code
func (c *ShortenlinkController) GetByCode(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid short code", nil)
		return
	}

	link, err := c.app.GetService().GetShortenlinkService().GetByCode(ctx.Request.Context(), code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(ctx, http.StatusNotFound, "Shortlink not found", nil)
			return
		}
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SendOne(ctx, gin.H{
		"code": link.ShortCode,
		"url":  link.OriginalURL,
		"id":   link.ID,
	}, http.StatusText(http.StatusOK), nil)
}

// PUT /api/shortenlinks/:code
func (c *ShortenlinkController) UpdateByCode(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid short code", nil)
		return
	}

	var req dto.UpdateShortenLinkDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	userIDVal, ok := ctx.Get("user_id")
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}
	userIDUUID, ok := userIDVal.(uuid.UUID)
	if !ok {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}

	updated, err := c.app.GetService().GetShortenlinkService().UpdateByCode(ctx.Request.Context(), code, req.OriginalURL, userIDUUID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(ctx, http.StatusNotFound, "Shortlink not found", nil)
			return
		}
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SendOne(ctx, gin.H{
		"code": updated.ShortCode,
		"url":  updated.OriginalURL,
		"id":   updated.ID,
	}, "Update success", nil)
}

//
// ==========================
// REDIRECT SHORTLINK (PUBLIC)
// ==========================
func (c *ShortenlinkController) Redirect(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid short code", nil)
		return
	}

	originalURL, err := c.app.GetService().
		GetShortenlinkService().
		Redirect(ctx.Request.Context(), code)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(ctx, http.StatusNotFound, "Shortlink not found", nil)
			return
		}
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, originalURL)
}
