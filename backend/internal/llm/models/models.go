package models

import "time"

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Model  string `json:"model"`
	Response  string    `json:"response"`
	CreatedAt time.Time `json:"created_at"`
	Done      bool      `json:"done"`
}

type GenTestRequest struct {
	Prompt string `json:"prompt"`
}