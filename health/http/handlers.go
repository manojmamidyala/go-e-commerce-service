package http

import (
	"mami/e-commerce/commons/logger"
	"mami/e-commerce/health/service"
	responsedto "mami/e-commerce/responseDto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	service service.IHealthService
}

func NewHealthHandler(service service.IHealthService) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

func (h *HealthHandler) Check(c *gin.Context) {
	status, err := h.service.Check(c)
	if err != nil {
		logger.Fatalf("Health check failed with error %v", err)
		responsedto.Error(c, http.StatusInternalServerError, err, "DOWN")
		return
	}

	responsedto.JSON(c, http.StatusOK, status)
}
