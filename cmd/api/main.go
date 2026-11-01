package main

import (
	"context"
	"fmt"
	"log"
	"service_blog_go/internal/config"
	"service_blog_go/internal/handler"
	"service_blog_go/internal/handler/middleware"
	"service_blog_go/internal/repository"
	"service_blog_go/internal/service"
	"service_blog_go/pkg/asynclock"
	"service_blog_go/pkg/database"
	"service_blog_go/pkg/opentelemetry"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	env := config.LoadEnvironment()
	port := env.Port

	// Init Otel
	ctx := context.Background()
	shutdown, err := opentelemetry.Init(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if shutdownErr := shutdown(shutdownCtx); shutdownErr != nil {
			log.Printf("Failed to shut down OpenTelemetry: %v", shutdownErr)
		}
	}()

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

	gin.DefaultWriter = log.Writer()
	gin.DefaultErrorWriter = log.Writer()
	route := gin.Default()
	route.Use(otelgin.Middleware(env.AppName))

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
