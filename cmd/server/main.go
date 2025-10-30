package main

import (
	"log"
	"os"

	"github.com/danielosbaldo/survey-app/internal/db"
	"github.com/danielosbaldo/survey-app/internal/models"
	"github.com/danielosbaldo/survey-app/internal/seed"
	"github.com/danielosbaldo/survey-app/internal/server"
	"github.com/gin-gonic/gin"
)

func main() {
	database, err := db.Open()
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	database.AutoMigrate(&models.Shop{}, &models.Ciudad{}, &models.Employee{}, &models.EmployeeShop{}, &models.Question{}, &models.Choice{}, &models.Response{})
	if err := seed.Run(database); err != nil {
		log.Fatalf("seed: %v", err)
	}

	s := server.New(database)
	r := s.Router()

	// ONLY trust local proxy (nginx on same box)
	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		log.Fatalf("trusted proxies: %v", err)
	}

	// production log level
	gin.SetMode(gin.ReleaseMode)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	s.Router().Run(":" + port)
}
