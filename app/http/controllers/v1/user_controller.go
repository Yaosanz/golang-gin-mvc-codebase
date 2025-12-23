package v1

import (
	"errors"
	"go-starter-app/app/http/dto"
	"go-starter-app/app/http/utils"
	"go-starter-app/helpers"
	"go-starter-app/interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	app interfaces.KernelDependencies
}

func NewUserController(app interfaces.KernelDependencies) *UserController {
	return &UserController{app: app}
}

// FindAll godoc
// @Summary List of users
// @Description Get a list of users with pagination
// @Tags User
// @Accept json
// @Produce json
// @Param paginate query bool false "turn pagination on/off" default(true)
// @Param page query int false "Page Number" default(1)
// @Param limit query int false "How many record per page" default(15)
// @Param sort_by query string false "sort list by field name" default(created_at)
// @Param sort_direction query string false "sort direction desc/asc" Enums(asc, desc) default(asc)
// @Param search query string false "search list by allowed fields"
// @Success 200 {object} swagger.ResponseOk{data=swagger.UserPaginatedResponse}
// @Failure 500 {object} swagger.ResponseError
// @Router /api/v1/users [get]
func (c *UserController) FindAll(ctx *gin.Context) {
	params := utils.ParseQueryParams(ctx, nil)
	params.Sanitize()
	// params.Filters["id"] = "2" // example injecting filter by id

	data, total, err := c.app.GetService().UserService.FindAll(ctx, params)
	if err != nil {
		utils.SendError(ctx, http.StatusNotFound, err.Error(), err)
		return
	}

	if params.Paginate {
		pagination := utils.NewPagination(
			helpers.ToInt(total),
			params.Page,
			params.Limit,
		)

		utils.SendPaginated(ctx, data, pagination, http.StatusText(http.StatusOK))
	} else {
		utils.SendList(ctx, data, http.StatusText(http.StatusOK))
	}
}

// FindByID godoc
// @Summary Get user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} swagger.ResponseOk{data=swagger.UserResponse}
// @Router /api/v1/users/{id} [get]
func (c *UserController) FindByID(ctx *gin.Context) {
	id, _ := utils.GetIDParam(ctx, "id")
	data, err := c.app.GetService().UserService.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(ctx, http.StatusNotFound, err.Error(), err)
			return
		}

		utils.SendError(ctx, http.StatusNotFound, err.Error(), err)
		return
	}

	utils.SendOne(ctx, data, http.StatusText(http.StatusOK), nil)
}

// Create godoc
// @Summary Create user
// @Tags User
// @Accept json
// @Produce json
// @Param payload body dto.CreateUserDTO true "Create user"
// @Success 200 {object} swagger.ResponseOk
// @Router /api/v1/users [post]
func (c *UserController) Create(ctx *gin.Context) {
	var req dto.CreateUserDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendError(ctx, http.StatusNotFound, "Invalid request data", err)
		return
	}

	// validate request
	validationErrors := c.app.GetValidator().ValidateStruct(&req)
	if validationErrors != nil {
		utils.SendError(ctx, http.StatusNotFound, "Validation error", validationErrors)
		return
	}

	err := c.app.GetService().UserService.Create(ctx, &req)
	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SendOne(ctx, nil, http.StatusText(http.StatusCreated), nil)
}

// Update godoc
// @Summary Update user
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param payload body dto.UpdateUserDTO true "Update user"
// @Success 200 {object} swagger.ResponseOk{data=swagger.UserResponse}
// @Router /api/v1/users/{id} [patch]
func (c *UserController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		utils.SendError(ctx, http.StatusInternalServerError, "Invalid ID", nil)
		return
	}

	// convert id to int64
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	var req dto.UpdateUserDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Invalid request data", err)
		return
	}

	// validate request
	validationErrors := c.app.GetValidator().ValidateStruct(&req)
	if validationErrors != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Validation error", validationErrors)
		return
	}

	err = c.app.GetService().UserService.Update(ctx, id, &req)
	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), err)
		return
	}

	utils.SendOne(ctx, nil, http.StatusText(http.StatusCreated), nil)
}

// Delete godoc
// @Summary Delete user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} swagger.ResponseOk
// @Failure 500 {object} swagger.ResponseError
// @Router /api/v1/users/{id} [delete]
func (c *UserController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	if idParam == "" {
		utils.SendError(ctx, http.StatusInternalServerError, "Invalid ID", nil)
		return
	}

	// convert id to int64
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.SendError(ctx, http.StatusInternalServerError, "Invalid ID format", err)
		return
	}

	err = c.app.GetService().UserService.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendError(ctx, http.StatusNotFound, "User not found", nil)
			return
		}

		utils.SendError(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SendOne(ctx, nil, http.StatusText(http.StatusOK), nil)
}
