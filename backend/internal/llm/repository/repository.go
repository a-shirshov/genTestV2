package repository

import (
	"encoding/json"
	"fmt"
	"gentest/internal/llm/models"
	"log"

	"github.com/go-resty/resty/v2"
)

type Repository struct {
	client *resty.Client
	llmConnURL string
}

func NewRepository(client *resty.Client, url string) *Repository {
	return &Repository{
		client: client,
		llmConnURL: url,
	}
}

func (r *Repository) GenerateResponse(ollamaReq *models.OllamaRequest) (*resty.Response, error) {
	jsonRequest, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, err
	}
	
	resp, err := r.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonRequest).
		SetDoNotParseResponse(true).
		Post(r.llmConnURL)

	if err != nil {
		return nil, fmt.Errorf("on generate response: %w", err)
	}

	log.Printf("Ollama response status: %s", resp.Status())

	return resp, nil
}