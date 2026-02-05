package controllers

import (
	"context"
	"go-starter-app/interfaces"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	db interfaces.KernelDependencies
}

func NewHealthController(deps interfaces.KernelDependencies) *HealthController {
	return &HealthController{
		db: deps,
	}
}

// GetHealth checks the health of the application and database
// @Summary Health Check
// @Description Check the health of the application and database connection
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (hc *HealthController) GetHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	dbErr := hc.checkDatabase(ctx)

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"database": map[string]interface{}{
			"status": "ok",
		},
	}

	statusCode := http.StatusOK

	if dbErr != nil {
		health["status"] = "degraded"
		health["database"].(map[string]interface{})["status"] = "error"
		health["database"].(map[string]interface{})["error"] = dbErr.Error()
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, health)
}

// checkDatabase verifies the database connection using GORM
func (hc *HealthController) checkDatabase(ctx context.Context) error {
	// Get database instance from app container
	appDeps, ok := hc.db.(interfaces.IAppDependencies)
	if !ok {
		return nil // If not available, skip check
	}

	db := appDeps.GetDB()
	if db == nil {
		return nil // Database not initialized, skip check
	}

	// Test database connectivity with a simple query
	var result string
	if err := db.WithContext(ctx).Raw("SELECT 'connection_ok'").Scan(&result).Error; err != nil {
		return err
	}

	return nil
}
