package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/enson89/sustainability-tracker-user-service/internal/config"
	"github.com/enson89/sustainability-tracker-user-service/internal/handler"
	"github.com/enson89/sustainability-tracker-user-service/internal/middleware"
	"github.com/enson89/sustainability-tracker-user-service/internal/repository"
	"github.com/enson89/sustainability-tracker-user-service/internal/service"
)

func main() {
	// Initialize zap logger.
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// Connect to PostgreSQL.
	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(authService)

	// Use gin.New() for a clean instance and add custom middleware.
	router := gin.New()
	// Recovery middleware for panics.
	router.Use(gin.Recovery())
	// Our centralized error-handling middleware.
	router.Use(middleware.ErrorHandler(logger))
	// Gin Logger middleware (optional, can be replaced by a custom zap-based logger).
	router.Use(gin.Logger())

	api := router.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.PUT("/profile", userHandler.UpdateProfile)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
