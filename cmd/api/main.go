package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"service_go_gin/internal/config"
	"service_go_gin/internal/handler"
	"service_go_gin/internal/handler/middleware"
	"service_go_gin/internal/repository"
	"service_go_gin/internal/service"
	"service_go_gin/pkg/asynclock"
	"service_go_gin/pkg/database"
	"service_go_gin/pkg/telemetry"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	env := config.LoadEnvironment()
	port := env.Port

	// OpenTelemetry (traces & metrics to SigNoz)
	shutdownTelemetry, err := telemetry.InitTelemetry(ctx, env)
	if err != nil {
		log.Fatalf("Failed init OpenTelemetry: %v", err)
	}

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
	// gin.Context.Value() fallback to Request.Context(), so span is propagated when *gin.Context passed as context.Context
	route.ContextWithFallback = true
	if env.OtelEnabled {
		route.Use(otelgin.Middleware(env.AppName))
	}

	v1 := route.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Ping)
		v1.POST("/auth", authHandler.Login)
		v1.GET("/user", middleware.AuthMiddleware(env), userHandler.GetUserByToken)
		v1.GET("/auth/refresh-token", middleware.RefreshTokenMiddleware(env), authHandler.RefreshToken)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: route,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed shutdown server: %v", err)
	}

	// Flush remaining spans & metrics before exit
	if err := shutdownTelemetry(shutdownCtx); err != nil {
		log.Printf("Failed shutdown OpenTelemetry: %v", err)
	}
}
