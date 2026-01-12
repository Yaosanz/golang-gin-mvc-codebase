package v1

import (
	"net/http"

	"go-starter-app/app/http/dto"
	"go-starter-app/app/http/utils"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	app interfaces.IAppDependencies
}

func NewAuthController(app interfaces.IAppDependencies) *AuthController {
	return &AuthController{
		app: app,
	}
}

// POST /api/v1/auth/login
func (ctl *AuthController) Login(c *gin.Context) {
	var req dto.LoginDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Validation error", err)
		return
	}

	token, err := ctl.app.GetService().
		AuthService.
		Login(c.Request.Context(), req.Username, req.Password)

	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, "Invalid username or password", nil)
		return
	}

	utils.SendSuccess(c, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
	}, "Login success")
}

// POST /api/v1/auth/register
func (ctl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Validation error", err)
		return
	}

	err := ctl.app.GetService().
		AuthService.
		Register(c.Request.Context(), req)

	if err != nil {
		utils.SendError(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.SendSuccess(c, nil, "Register success")
}