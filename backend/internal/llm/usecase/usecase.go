package usecase

import (
	"gentest/internal/llm"
	"gentest/internal/llm/models"

	"github.com/go-resty/resty/v2"
)

type Usecase struct {
	repo llm.Repository
}

func NewUsecase(repo llm.Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) GenerateResponse(ollamaReq *models.OllamaRequest) (*resty.Response, error) {
	return u.repo.GenerateResponse(ollamaReq)
}