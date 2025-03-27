package main

import (
	"log"
	"os"

	"github.com/enson89/sustainability-tracker-user-service/internal/config"
	"github.com/enson89/sustainability-tracker-user-service/internal/handler"
	"github.com/enson89/sustainability-tracker-user-service/internal/repository"
	"github.com/enson89/sustainability-tracker-user-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to PostgreSQL
	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(authService)

	router := gin.Default()
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
