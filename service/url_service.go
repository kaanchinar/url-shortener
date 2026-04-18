package service

import (
	"context"
	"strconv"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/kaanchinar/url-shortener/dto"
	"github.com/kaanchinar/url-shortener/model"
)

type URLRepository interface {
	CreateUrl(ctx context.Context, url model.URL) error
	GetUrlById(ctx context.Context, id string) (*model.URL, error)
}

type URLService struct {
	repo URLRepository
}

func (s *URLService) ShortenUrl(ctx context.Context, req dto.CreateShortURLRequest) (string, error) {
	url := model.URL{
		ID:          generateUniqueID(req.URL), // Implement this function to generate a unique ID
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

func generateUniqueID(originalURL string) string {
	return strconv.FormatUint(xxhash.Sum64String(originalURL), 10)
}
