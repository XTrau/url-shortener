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

func (repo *RedirectStatisticsPostgresRepository) RedirrectCount(ctx context.Context, urlID int, start time.Time, end time.Time) (int, error) {
	query := `SELECT COUNT(1) FROM redirects WHERE url_id=$1 AND ($2 <= redirected_at AND redirected_at < $3)`
	row := repo.db.QueryRowContext(ctx, query, urlID, start, end)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RedirectStatisticsPostgresRepository) IpCount(ctx context.Context, urlID int, start time.Time, end time.Time) (int, error) {
	query := `SELECT COUNT(1) FROM (SELECT DISTINCT ip FROM redirects WHERE url_id=$1 AND ($2 <= redirected_at AND redirected_at < $3))`
	row := repo.db.QueryRowContext(ctx, query, urlID, start, end)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RedirectStatisticsPostgresRepository) TotalRedirrectCount(ctx context.Context, urlID int) (int, error) {
	query := `SELECT COUNT(1) FROM redirects WHERE url_id=$1`
	row := repo.db.QueryRowContext(ctx, query, urlID)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RedirectStatisticsPostgresRepository) TotalIpCount(ctx context.Context, urlID int) (int, error) {
	query := `SELECT COUNT(1) FROM (SELECT DISTINCT ip FROM redirects WHERE url_id=$1)`
	row := repo.db.QueryRowContext(ctx, query, urlID)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
