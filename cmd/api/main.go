package main

import (
	"fmt"
	"log"
	"service_go_gin/internal/config"
	"service_go_gin/internal/handler"
	"service_go_gin/internal/handler/middleware"
	"service_go_gin/internal/repository"
	"service_go_gin/internal/service"
	"service_go_gin/pkg/asynclock"
	"service_go_gin/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	env := config.LoadEnvironment()
	port := env.Port

	// connection db postgres
	db_postgres := database.InitPostgres(env)
	asyncLock := asynclock.NewSycnLock()

	healthService := service.NewHealthService(env)
	healthHandler := handler.NewHealthHandler(healthService)

	userRepo := repository.NewUserRepository(db_postgres)
	userService := service.NewUserService(userRepo, env)
	userHandler := handler.NewUserHandler(userService)

	authRepo := repository.NewauthTokenRepository(db_postgres)
	authService := service.NewAuthService(&service.DependsAuthService{
		UserRepo: userRepo,
		AuthRepo: authRepo,
		Env:      env,
		SyncLock: asyncLock,
		DB:       db_postgres,
	})
	authHandler := handler.NewAuthHandler(authService, env)

	route := gin.Default()

	v1 := route.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Ping)
		v1.POST("/auth", authHandler.Login)
		v1.GET("/user", middleware.AuthMiddleware(env), userHandler.GetUserByToken)
		v1.GET("/auth/refresh-token", middleware.RefreshTokenMiddleware(env), authHandler.RefreshToken)
	}

	serverAddr := fmt.Sprintf(":%s", port)
	if err := route.Run(serverAddr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
