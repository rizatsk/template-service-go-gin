package service

import (
	"context"
	"service_go_gin/internal/config"
	"service_go_gin/internal/domain"
)

type healthService struct {
	env *config.Environment
}

func NewHealthService(env *config.Environment) domain.HealthService {
	return &healthService{
		env: env,
	}
}

func (service *healthService) CheckHealth(ctx context.Context) *domain.DataHealthResponse {
	return &domain.DataHealthResponse{
		AppName: service.env.AppName,
		AppEnv:  service.env.AppEnv,
	}
}
