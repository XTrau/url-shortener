package repositories

import (
	"context"
	"database/sql"
	"time"
)

type RedirectStatisticsPostgresRepository struct {
	db *sql.DB
}

func NewRedirectStatisticsPostgresRepository(db *sql.DB) *RedirectStatisticsPostgresRepository {
	return &RedirectStatisticsPostgresRepository{db}
}

func (repo *RedirectStatisticsPostgresRepository) RedirrectCount(ctx context.Context, urlID int, t *time.Time) (int, error) {
	var row *sql.Row

	if t == nil {
		query := `SELECT COUNT(1) FROM redirects WHERE url_id=$1`
		row = repo.db.QueryRowContext(ctx, query, urlID)
	} else {
		query := `SELECT COUNT(1) FROM redirects WHERE url_id=$1 AND redirected_at>=$2`
		row = repo.db.QueryRowContext(ctx, query, urlID, *t)
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RedirectStatisticsPostgresRepository) IpCount(ctx context.Context, urlID int, t *time.Time) (int, error) {
	var row *sql.Row
	if t == nil {
		query := `SELECT COUNT(1) FROM (SELECT DISTINCT ip FROM redirects WHERE url_id=$1)`
		row = repo.db.QueryRowContext(ctx, query, urlID)
	} else {
		query := `SELECT COUNT(1) FROM (SELECT DISTINCT ip FROM redirects WHERE url_id=$1 AND redirected_at>=$2)`
		row = repo.db.QueryRowContext(ctx, query, urlID, *t)
	}

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
