package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kaanchinar/url-shortener/model"
)

type DB interface {
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type URLRepository struct {
	db DB
}

func NewURLRepository(db DB) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) CreateUrl(ctx context.Context, url model.URL) error {
	query := `INSERT INTO urls (id, original_url, created_at, updated_at) VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(ctx, query, url.ID, url.OriginalURL, url.CreatedAt, url.UpdatedAt)
	return err

}

func (r *URLRepository) GetUrlById(ctx context.Context, id string) (*model.URL, error) {
	query := `SELECT id, original_url, created_at, updated_at FROM urls WHERE id = $1`
	var url model.URL
	err := r.db.QueryRow(ctx, query, id).Scan(&url.ID, &url.OriginalURL, &url.CreatedAt, &url.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &url, nil
}
