package service

import (
	"service_go_gin/internal/config"
	"service_go_gin/internal/domain"

	"github.com/gin-gonic/gin"
)

type UserService struct {
	userRepo domain.UserRepository
	env      *config.Environment
}

func NewUserService(userRepo domain.UserRepository, env *config.Environment) domain.UserService {
	return &UserService{
		userRepo: userRepo,
		env:      env,
	}
}

func (userService UserService) GetUser(ctx *gin.Context) *domain.GetUser {
	return &domain.GetUser{
		ID:        ctx.MustGet("user_id").(string),
		Name:      ctx.MustGet("name").(string),
		Email:     ctx.MustGet("email").(string),
		Role:      ctx.MustGet("role").(string),
		CreatedAt: ctx.MustGet("created_at").(string),
		UpdatedAt: ctx.MustGet("updated_at").(string),
	}
}
