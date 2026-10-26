package handler

import (
	"net/http"
	"service_go_gin/internal/domain"
	"service_go_gin/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (user *UserHandler) GetUserByToken(ctx *gin.Context) {
	dataUser := user.userService.GetUser(ctx)
	response.Success(
		ctx,
		http.StatusOK,
		"Success get data user",
		dataUser,
	)
}
