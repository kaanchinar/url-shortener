package service

import (
	"context"
	"time"

	"github.com/kaanchinar/url-shortener/dto"
	"github.com/kaanchinar/url-shortener/model"
	"github.com/kaanchinar/url-shortener/utils"
)

type URLRepository interface {
	CreateUrl(ctx context.Context, url model.URL) error
	GetUrlById(ctx context.Context, id string) (*model.URL, error)
}

type URLService struct {
	repo URLRepository
}

func NewURLService(repo URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (s *URLService) ShortenUrl(ctx context.Context, req dto.CreateShortURLRequest) (string, error) {
	url := model.URL{
		ID:          utils.GenerateUniqueID(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		OriginalURL: req.URL,
	}

	err := s.repo.CreateUrl(ctx, url)
	if err != nil {
		return "there was an error on shorten url", err
	}

	return url.ID, nil
}

func (s *URLService) GetUrlById(ctx context.Context, id string) (*model.URL, error) {
	url, err := s.repo.GetUrlById(ctx, id)
	if err != nil {
		return nil, err
	}
	return url, nil
}
