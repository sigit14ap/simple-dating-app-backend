package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sigit14ap/simple-dating-app-backend/internal/databases"
	handlers "github.com/sigit14ap/simple-dating-app-backend/internal/handlers"
	"github.com/sigit14ap/simple-dating-app-backend/internal/repositories"
	"github.com/sigit14ap/simple-dating-app-backend/internal/services"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Application start ...")
	log.Info("Logger initialized ...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := databases.DatabaseInitialize()
	if err != nil {
		log.Fatalf("Error initialize database: %s", err)
	}

	userRepo := repositories.NewUserRepository(db)
	matchRepo := repositories.NewMatchRepository(db)

	userService := services.NewUserService(userRepo, matchRepo)
	matchService := services.NewMatchService(matchRepo, userRepo)

	handler := handlers.NewHandler(userService, matchService)

	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/signup", handler.Signup)
		api.POST("/login", handler.Login)
		api.GET("/next-user", handler.GetNextUser)
		api.POST("/swipe", handler.Swipe)
		api.POST("/buy-package-unlimited-swipe", handler.BuyPackageUnlimitedSwipe)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
