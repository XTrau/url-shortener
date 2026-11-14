// Репозиторий для работы с url
package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"urlshortener/internal/apperrors"
	"urlshortener/internal/domain"
)

type UrlPostgresRepository struct {
	db *sql.DB
}

func NewUrlPostgresRepository(db *sql.DB) *UrlPostgresRepository {
	return &UrlPostgresRepository{db}
}

func (repo *UrlPostgresRepository) CreateSlugIfNotExists(ctx context.Context, slugText string) (domain.Slug, error) {
	slog.Debug("Creating if not exists slug", slog.String("slug", slugText))

	query := "INSERT INTO slugs(slug) VALUES ($1) ON CONFLICT (slug) DO NOTHING RETURNING id, next_number"
	row := repo.db.QueryRowContext(ctx, query, slugText)

	var id, nextNumber int

	if err := row.Scan(&id, &nextNumber); err != nil {
		return domain.Slug{}, err
	}

	return domain.Slug{
		ID:         id,
		Text:       slugText,
		NextNumber: nextNumber,
	}, nil
}

func (repo *UrlPostgresRepository) CreateUrl(ctx context.Context, urlText string, slug domain.Slug) (domain.Url, error) {
	slog.Debug("Inserting url to database", slog.String("url", urlText), slog.String("slug", slug.Text))

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Url{}, err
	}
	defer tx.Rollback()

	query := "UPDATE slugs SET next_number = next_number + 1 WHERE id=$1 RETURNING next_number - 1"
	row := tx.QueryRowContext(ctx, query, slug.ID)

	var slugNumber int
	if err := row.Scan(&slugNumber); err != nil {
		return domain.Url{}, err
	}

	query = "INSERT INTO urls (url, slug_id, slug_number) VALUES ($1, $2, $3) RETURNING id"
	row = repo.db.QueryRowContext(ctx, query, urlText, slug.ID, slugNumber)

	var id int
	if err := row.Scan(&id); err != nil {
		return domain.Url{}, err
	}

	if err := tx.Commit(); err != nil {
		return domain.Url{}, err
	}

	return domain.Url{
		ID:  id,
		Url: urlText,
		Slug: &domain.UrlSlug{
			Text:   slug.Text,
			Number: slugNumber,
		},
	}, nil
}

func (repo *UrlPostgresRepository) GetUrlBySlug(ctx context.Context, slug domain.UrlSlug) (domain.Url, error) {
	slog.Debug("Getting url from database", slog.String("slugText", slug.Text), slog.Int("slugNumber", slug.Number))

	query := `
	SELECT u.id, u.url, s.slug, u.slug_number
	FROM urls u
	JOIN slugs s ON u.slug_id=s.id
	WHERE slug_number=$1 AND s.slug=$2
	`
	row := repo.db.QueryRowContext(ctx, query, slug.Number, slug.Text)

	var url string
	err := row.Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Url{}, apperrors.ErrUrlNotFound
	}
	return domain.Url{Url: url}, err
}

func (repo *UrlPostgresRepository) GetSlugByUrl(ctx context.Context, urlText string) (domain.UrlSlug, error) {
	slog.Debug("Getting url from database", slog.String("url", urlText))

	query := `
	SELECT s.slug, u.slug_number
	FROM urls u
	JOIN slugs s ON u.slug_id=s.id
	WHERE u.url=$1
	`
	row := repo.db.QueryRowContext(ctx, query, urlText)

	var slug string
	var slugNumber int

	err := row.Scan(&slug, &slugNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.UrlSlug{}, apperrors.ErrSlugNotFound
	}

	return domain.UrlSlug{
		Text:   slug,
		Number: slugNumber,
	}, err
}
