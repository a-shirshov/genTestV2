package llm

import (
	"gentest/internal/llm/models"

	"github.com/go-resty/resty/v2"
)

type Usecase interface {
	GenerateResponse(ollamaReq *models.OllamaRequest) (*resty.Response, error)
}