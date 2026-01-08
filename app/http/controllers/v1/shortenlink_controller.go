package v1

import (
	"errors"
	"net/http"

	"go-starter-app/app/http/dto"
	"go-starter-app/app/http/utils"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
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
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}

	data, err := c.app.GetService().
		ShortenlinkService.
		Create(ctx.Request.Context(), req.OriginalURL, userID)

	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SendOne(ctx, data, http.StatusText(http.StatusCreated), nil)
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
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}

	data, err := c.app.GetService().
		ShortenlinkService.
		FindAll(ctx.Request.Context(), userID)
	if err != nil {
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
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}

	data, err := c.app.GetService().
		ShortenlinkService.
		FindByID(ctx.Request.Context(), id, userID)
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
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}

	data, err := c.app.GetService().
		ShortenlinkService.
		Update(ctx.Request.Context(), id, req.OriginalURL, userID)
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
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		utils.SendError(ctx, http.StatusUnauthorized, "Invalid user id", nil)
		return
	}

	err := c.app.GetService().
		ShortenlinkService.
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

	data, err := c.app.GetService().
		ShortenlinkService.
		GetByCode(ctx.Request.Context(), code)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(ctx, http.StatusNotFound, "Shortlink not found", nil)
			return
		}
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, data.OriginalURL)
}
