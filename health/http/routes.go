package http

import (
	"mami/e-commerce/config"
	"mami/e-commerce/health/service"

	"github.com/gin-gonic/gin"
)

func Routes(r *gin.RouterGroup, sqlDB config.IDatabase) {
	healthService := service.NewHealthService()
	healthHandler := NewHealthHandler(healthService)

	healthRoute := r.Group("/health")
	{
		healthRoute.GET("/check", healthHandler.Check)
	}
}
