package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kaanchinar/url-shortener/model"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) CreateUrl(ctx context.Context, url model.URL) error {
	query := `INSERT INTO urls (id, original_url, created_at, updated_at) VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query, url.ID, url.OriginalURL, url.CreatedAt, url.UpdatedAt)
	return err

}

func (r *URLRepository) GetUrlById(ctx context.Context, id string) (*model.URL, error) {
	query := `SELECT id, original_url, created_at, updated_at FROM urls WHERE id = $1`
	var url model.URL
	err := r.db.QueryRowContext(ctx, query, id).Scan(&url.ID, &url.OriginalURL, &url.CreatedAt, &url.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &url, nil
}
