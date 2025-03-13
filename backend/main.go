package main

import (
	"gentest/internal/llm/delivery"
	"gentest/internal/llm/repository"
	"gentest/internal/llm/usecase"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type Body struct {
	Prompt string `json:"prompt"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func main() {
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins, or specify your allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	client := resty.New()
	repo := repository.NewRepository(client, "http://localhost:11434/api/generate")
	uc := usecase.NewUsecase(repo)
	delivery := delivery.NewDelivery(uc)

	r.POST("/ask", delivery.GenerateFirstTest)
	r.Run(":3000")
}

