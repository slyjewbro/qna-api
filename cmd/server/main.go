package main

import (
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"qna-api/internal/config"
	"qna-api/internal/handler"
	"qna-api/internal/model"
	"qna-api/internal/repository"
	"qna-api/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := gorm.Open(postgres.Open(cfg.GetDBConnectionString()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate models
	if err := db.AutoMigrate(&model.Question{}, &model.Answer{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize layers
	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	// Setup routes
	router := h.InitRoutes()

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}
