package repositories

import (
	"context"
	"database/sql"
	"urlshortener/internal/domain"
)

type RedirectsPostgresRepository struct {
	db *sql.DB
}

func NewRedirectsPostgresRepository(db *sql.DB) *RedirectsPostgresRepository {
	return &RedirectsPostgresRepository{db}
}

func (repo *RedirectsPostgresRepository) Save(ctx context.Context, redirect domain.RedirectInfo) error {
	query := `INSERT INTO redirects
	(url_id, redirected_at, ip, user_agent)
	VALUES
	($1, $2, $3, $4)
	`

	return repo.db.QueryRowContext(ctx, query, redirect.UrlID, redirect.CreatedAt, redirect.IP, redirect.UserAgent).Err()
}
