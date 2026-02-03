package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-starter-app/app/http/dto"
	"go-starter-app/app/repositories"
	"go-starter-app/app/services"
	"go-starter-app/app/validation"
	"go-starter-app/config"
	"go-starter-app/interfaces"
	"go-starter-app/pkg/google"
	"go-starter-app/pkg/scheduler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// --- Mocks ---

type mockAuthService struct{
    lastRefreshInput string
}

func (m *mockAuthService) SecureLogin(ctx context.Context, username, password string) (string, string, interface{}, error) {
    return "", "", nil, nil
}

func (m *mockAuthService) Refresh(ctx context.Context, refreshToken string) (string, string, interface{}, error) {
    m.lastRefreshInput = refreshToken
    lc := &services.LoginContext{
        UserID:   uuid.New(),
        Username: "tester",
        Roles:    []string{"user"},
    }
    return "new_access", "new_refresh", lc, nil
}

func (m *mockAuthService) Register(ctx context.Context, r *dto.RegisterDTO) error { return nil }
func (m *mockAuthService) CheckPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) { return true, nil }
func (m *mockAuthService) CheckRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error) { return true, nil }
func (m *mockAuthService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) { return nil, nil }
func (m *mockAuthService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) { return nil, nil }
func (m *mockAuthService) IsUserActive(ctx context.Context, userID uuid.UUID) (bool, error) { return true, nil }
func (m *mockAuthService) RequirePermission(permission string) func(ctx context.Context, userID uuid.UUID) error { return nil }
func (m *mockAuthService) RequireRole(role string) func(ctx context.Context, userID uuid.UUID) error { return nil }
func (m *mockAuthService) RequireActiveUser() func(ctx context.Context, userID uuid.UUID) error { return nil }
func (m *mockAuthService) InvalidateUserSession(ctx context.Context, sessionID string) error { return nil }
func (m *mockAuthService) IsSessionValid(ctx context.Context, sessionID string) (bool, error) { return true, nil }
func (m *mockAuthService) InvalidateUserAuthCache(ctx context.Context, username string) error { return nil }

type mockServiceContainer struct{
    auth interfaces.IAuthService
}

func (m *mockServiceContainer) GetAuthService() interfaces.IAuthService { return m.auth }
func (m *mockServiceContainer) GetUserService() interfaces.IUserService { return nil }
func (m *mockServiceContainer) GetShortenlinkService() interfaces.IShortenlinkService { return nil }
func (m *mockServiceContainer) GetPermissionService() interfaces.IPermissionService { return nil }

type mockAppDeps struct{
    svc interfaces.IServiceContainer
}

func (m *mockAppDeps) GetConfig() *config.Config { return nil }
func (m *mockAppDeps) GetDB() *gorm.DB { return nil }
func (m *mockAppDeps) GetDBTransaction() *gorm.DB { return nil }
func (m *mockAppDeps) GetDBWithContext(ctx context.Context) *gorm.DB { return nil }
func (m *mockAppDeps) GetRepo() *repositories.RepoContainer { return nil }
func (m *mockAppDeps) GetService() interfaces.IServiceContainer { return m.svc }
func (m *mockAppDeps) GetValidator() *validation.AppValidator { return nil }
func (m *mockAppDeps) GetScheduler() *scheduler.Scheduler { return nil }
func (m *mockAppDeps) GetMinio() *minio.Client { return nil }
func (m *mockAppDeps) GetFcm() *google.FCM { return nil }
func (m *mockAppDeps) GetRedis() *redis.Client { return nil }

// --- Helper to perform request ---

func performRefresh(t *testing.T, router *gin.Engine, req *http.Request) (int, map[string]interface{}) {
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        return w.Code, nil
    }
    var resp struct{
        Data map[string]interface{} `json:"data"`
    }
    if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
        t.Fatalf("failed to parse response: %v", err)
    }
    return w.Code, resp.Data
}

func TestRefreshTokenIntakePaths(t *testing.T) {
    gin.SetMode(gin.TestMode)
    mockAuth := &mockAuthService{}
    mockSvc := &mockServiceContainer{auth: mockAuth}
    app := &mockAppDeps{svc: mockSvc}
    ctl := NewAuthController(app)

    r := gin.New()
    r.POST("/api/auth/refresh", ctl.Refresh)

    // 1) Authorization: Bearer
    token := "tok_bearer_123"
    req1, _ := http.NewRequest("POST", "/api/auth/refresh", nil)
    req1.Header.Set("Authorization", "Bearer "+token)
    code, data := performRefresh(t, r, req1)
    if code != http.StatusOK || data["token"] == nil || data["refresh_token"] == nil || mockAuth.lastRefreshInput != token {
        t.Fatalf("bearer auth failed: code=%d, last=%s", code, mockAuth.lastRefreshInput)
    }

    // 2) X-Refresh-Token (raw)
    token2 := "tok_alt_header"
    req2, _ := http.NewRequest("POST", "/api/auth/refresh", nil)
    req2.Header.Set("X-Refresh-Token", token2)
    code, _ = performRefresh(t, r, req2)
    if code != http.StatusOK || mockAuth.lastRefreshInput != token2 {
        t.Fatalf("alt header failed: code=%d, last=%s", code, mockAuth.lastRefreshInput)
    }

    // 3) Refresh-Token: Bearer <token>
    token3 := "tok_alt_bearer"
    req3, _ := http.NewRequest("POST", "/api/auth/refresh", nil)
    req3.Header.Set("Refresh-Token", "Bearer "+token3)
    code, _ = performRefresh(t, r, req3)
    if code != http.StatusOK || mockAuth.lastRefreshInput != token3 {
        t.Fatalf("alt header bearer failed: code=%d, last=%s", code, mockAuth.lastRefreshInput)
    }

    // 4) Query param
    token4 := "tok_query"
    req4, _ := http.NewRequest("POST", "/api/auth/refresh?refresh_token="+token4, nil)
    code, _ = performRefresh(t, r, req4)
    if code != http.StatusOK || mockAuth.lastRefreshInput != token4 {
        t.Fatalf("query param failed: code=%d, last=%s", code, mockAuth.lastRefreshInput)
    }

    // 5) JSON body
    token5 := "tok_json"
    body := `{"refresh_token":"` + token5 + `"}`
    req5, _ := http.NewRequest("POST", "/api/auth/refresh", strings.NewReader(body))
    req5.Header.Set("Content-Type", "application/json")
    code, _ = performRefresh(t, r, req5)
    if code != http.StatusOK || mockAuth.lastRefreshInput != token5 {
        t.Fatalf("json body failed: code=%d, last=%s", code, mockAuth.lastRefreshInput)
    }
}
