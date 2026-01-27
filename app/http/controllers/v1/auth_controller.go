package v1

import (
	"errors"
	"net/http"
	"strings"

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

	accessToken, refreshToken, loginCtx, err := ctl.app.GetService().
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
		"token":          accessToken,
		"refresh_token":  refreshToken,
		"token_type":     "Bearer",
		"user": gin.H{
			"id":       loginData.UserID.String(),
			"username": loginData.Username,
			"email":    userEmail,
			"roles":    loginData.Roles,
		},
	}

	utils.SendOne(c, data, "Login success", nil)
}

// POST /api/auth/refresh
func (ctl *AuthController) Refresh(c *gin.Context) {
	// Accept refresh token via Authorization: Bearer <token>, or fallback to JSON body
	var refreshTokenInput string

	// 1) Try Authorization header first
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		authHeader = c.Request.Header.Get("Authorization")
	}
	if authHeader != "" {
		parts := strings.Fields(authHeader)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			refreshTokenInput = strings.TrimSpace(parts[1])
		}
	}

	// 1a) Try alternative headers commonly used when proxies strip Authorization
	if refreshTokenInput == "" {
		altHeaders := []string{"X-Refresh-Token", "Refresh-Token", "X-Authorization"}
		for _, h := range altHeaders {
			val := strings.TrimSpace(c.GetHeader(h))
			if val == "" {
				val = strings.TrimSpace(c.Request.Header.Get(h))
			}
			if val != "" {
				// If header looks like "Bearer <token>", parse it; otherwise assume it's the token
				parts := strings.Fields(val)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					refreshTokenInput = strings.TrimSpace(parts[1])
				} else {
					refreshTokenInput = val
				}
				break
			}
		}
	}

	// 1a-extended) Iterate all headers to find any value containing "Bearer"
	if refreshTokenInput == "" {
		for name, vals := range c.Request.Header {
			// Check common mis-cased variants of Authorization
			if strings.EqualFold(name, "Authorization") && len(vals) > 0 {
				parts := strings.Fields(vals[0])
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					refreshTokenInput = strings.TrimSpace(parts[1])
					break
				}
			}
			// Fallback: search any header value for Bearer scheme
			for _, v := range vals {
				if strings.HasPrefix(strings.TrimSpace(v), "Bearer ") {
					refreshTokenInput = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(v), "Bearer "))
					break
				}
			}
			if refreshTokenInput != "" {
				break
			}
		}
	}

	// 1b) Try query params as a last header-less option
	if refreshTokenInput == "" {
		qp := strings.TrimSpace(c.Query("refresh_token"))
		if qp == "" {
			qp = strings.TrimSpace(c.Query("token"))
		}
		if qp != "" {
			refreshTokenInput = qp
		}
	}

	// 2) Fallback to JSON body only if Content-Type indicates JSON
	if refreshTokenInput == "" {
		ct := strings.ToLower(strings.TrimSpace(c.GetHeader("Content-Type")))
		if strings.Contains(ct, "application/json") {
			var req dto.RefreshTokenDTO
			if err := c.ShouldBindJSON(&req); err == nil {
				refreshTokenInput = req.RefreshToken
			}
		}
	}

	if refreshTokenInput == "" {
		utils.SendError(c, http.StatusBadRequest, "Refresh token required", nil)
		return
	}

	accessToken, refreshToken, loginCtx, err := ctl.app.GetService().
		GetAuthService().
		Refresh(c.Request.Context(), refreshTokenInput)

	if err != nil {
		utils.SendError(c, http.StatusUnauthorized, "Invalid or expired refresh token", nil)
		return
	}

	loginData, _ := loginCtx.(*services.LoginContext)
	if loginData == nil {
		utils.SendError(c, http.StatusInternalServerError, "Login context missing", nil)
		return
	}

	data := gin.H{
		"token":          accessToken,
		"refresh_token":  refreshToken,
		"token_type":     "Bearer",
		"user": gin.H{
			"id":       loginData.UserID.String(),
			"username": loginData.Username,
			"roles":    loginData.Roles,
		},
	}

	utils.SendOne(c, data, "Token refresh success", nil)
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
