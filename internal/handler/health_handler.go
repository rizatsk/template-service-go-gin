package handler

import (
	"net/http"
	"service_go_gin/internal/domain"
	"service_go_gin/pkg/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	healthService domain.HealthService
}

func NewHealthHandler(service domain.HealthService) *HealthHandler {
	return &HealthHandler{healthService: service}
}

func (health *HealthHandler) Ping(ctx *gin.Context) {
	reponseHealth := health.healthService.CheckHealth(ctx.Request.Context())

	response.Success(
		ctx,
		http.StatusOK,
		"Service is healthy",
		reponseHealth,
	)
}
