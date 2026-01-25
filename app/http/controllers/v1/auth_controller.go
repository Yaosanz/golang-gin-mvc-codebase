package v1

import (
	"errors"
	"net/http"

	"go-starter-app/app/http/dto"
	"go-starter-app/app/http/utils"
	"go-starter-app/app/services"
	"go-starter-app/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	app interfaces.IAppDependencies
}

func NewAuthController(app interfaces.IAppDependencies) *AuthController {
	return &AuthController{
		app: app,
	}
}

// POST /api/auth/login
func (ctl *AuthController) Login(c *gin.Context) {
	var req dto.LoginDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Validation error", err)
		return
	}

	token, loginCtx, err := ctl.app.GetService().
		GetAuthService().
		SecureLogin(c.Request.Context(), req.Username, req.Password)

	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, "Invalid username or password", nil)
		return
	}

	loginData, _ := loginCtx.(*services.LoginContext)
	if loginData == nil {
		utils.SendError(c, http.StatusInternalServerError, "Login context missing", nil)
		return
	}
	userSvc := ctl.app.GetService().GetUserService()
	userEmail := ""
	if user, _ := userSvc.FindById(c.Request.Context(), loginData.UserID.String()); user != nil {
		userEmail = user.Email
	}

	data := gin.H{
		"token":      token,
		"token_type": "Bearer",
		"user": gin.H{
			"id":       loginData.UserID.String(),
			"username": loginData.Username,
			"email":    userEmail,
			"roles":    loginData.Roles,
		},
	}

	utils.SendOne(c, data, "Login success", nil)
}

// POST /api/auth/register
func (ctl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Validation error", err)
		return
	}

	err := ctl.app.GetService().
		GetAuthService().
		Register(c.Request.Context(), &req)

	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, services.ErrUsernameExists) || errors.Is(err, services.ErrEmailExists) {
			status = http.StatusConflict
		}
		utils.SendError(c, status, err.Error(), err.Error())
		return
	}

	user, _ := ctl.app.GetService().GetUserService().FindByUsername(c.Request.Context(), req.Username)
	status := http.StatusCreated
	resp := gin.H{}
	if user != nil {
		resp = gin.H{
			"id":       user.ID.String(),
			"username": user.Username,
			"email":    user.Email,
		}
	}
	utils.SendOne(c, resp, "Register success", &status)
}

// GET /api/auth/profile
func (ctl *AuthController) Profile(c *gin.Context) {
	userIDVal, ok := c.Get("user_id")
	if !ok {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	userID := userIDVal.(uuid.UUID)
	user, err := ctl.app.GetService().GetUserService().FindById(c.Request.Context(), userID.String())
	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, "User not found", err)
		return
	}

	authSvc := ctl.app.GetService().GetAuthService().(*services.SecureAuthService)
	roles, _ := authSvc.GetUserRoles(c.Request.Context(), userID)
	perms, _ := authSvc.GetUserPermissions(c.Request.Context(), userID)

	utils.SendOne(c, gin.H{
		"id":          user.ID.String(),
		"username":    user.Username,
		"email":       user.Email,
		"roles":       roles,
		"permissions": perms,
	}, "Profile", nil)
}
