package handler

import (
	"net/http"
	"service_go_gin/internal/config"
	"service_go_gin/internal/domain"
	"service_go_gin/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService domain.AuthService
	env         *config.Environment
}

func NewAuthHandler(authService domain.AuthService, env *config.Environment) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		env:         env,
	}
}

func (authHandler *AuthHandler) Login(ctx *gin.Context) {
	var req domain.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		formatedErr := response.FormValidationError(err)
		response.Failed(
			ctx,
			response.NewAppError(400, "Bad request", formatedErr),
		)

		return
	}

	token, err := authHandler.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		response.Failed(
			ctx,
			err,
		)
		return
	}

	ctx.SetCookie(
		"token",           // name cookie
		token.TokenString, // value cookie
		authHandler.env.JWTExpiredHours*60*60, // expired from second
		"/",   // appicable path
		"",    // domain, (for localhost is blank)
		false, // secure  (set true when using ssl)
		true,  // http only (javascript cannot read this cookie)
	)

	ctx.SetCookie(
		"refreshToken",           // name cookie
		token.RefreshTokenString, // value cookie
		authHandler.env.JWTRfreshTokenExpiredDays*24*60*int(time.Hour), // expired from second
		"/",   // appicable path
		"",    // domain, (for localhost is blank)
		false, // secure  (set true when using ssl)
		true,  // http only (javascript cannot read this cookie)
	)

	response.Success(
		ctx,
		http.StatusOK,
		"Login success",
		nil,
	)
}

func (authHandler *AuthHandler) RefreshToken(ctx *gin.Context) {
	token, err := authHandler.authService.RefreshToken(ctx)
	if err != nil {
		response.Failed(
			ctx,
			err,
		)
		return
	}

	ctx.SetCookie(
		"token",           // name cookie
		token.TokenString, // value cookie
		authHandler.env.JWTExpiredHours*60*60, // expired from second
		"/",   // appicable path
		"",    // domain, (for localhost is blank)
		false, // secure  (set true when using ssl)
		true,  // http only (javascript cannot read this cookie)
	)

	ctx.SetCookie(
		"refreshToken",           // name cookie
		token.RefreshTokenString, // value cookie
		authHandler.env.JWTRfreshTokenExpiredDays*24*60*int(time.Hour), // expired from second
		"/",   // appicable path
		"",    // domain, (for localhost is blank)
		false, // secure  (set true when using ssl)
		true,  // http only (javascript cannot read this cookie)
	)

	response.Success(
		ctx,
		http.StatusOK,
		"Refresh token success",
		nil,
	)
}
