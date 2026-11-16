package usecases

import (
	"context"
	"fmt"
	"time"
	"urlshortener/internal/domain"
)

type StatsGetter interface {
	TotalRedirrectCount(ctx context.Context, urlID int) (int, error)
	TotalIpCount(ctx context.Context, urlID int) (int, error)
	// Возвращает количество редиректов в интервале [start, end).
	RedirrectCount(ctx context.Context, urlID int, start time.Time, end time.Time) (int, error)
	// Возвращает количество уникальных ip в интервале [start, end).
	IpCount(ctx context.Context, urlID int, start time.Time, end time.Time) (int, error)
}

type StatisticsUseCases struct {
	sg StatsGetter
}

func NewStatisticsUseCases(sg StatsGetter) *StatisticsUseCases {
	return &StatisticsUseCases{sg}
}

func (uc *StatisticsUseCases) Summary(ctx context.Context, urlID int) (domain.SummaryStatistics, error) {
	now := time.Now()
	now24h := now.Add(-time.Hour * 24)
	now7d := now.Add(-time.Hour * 24 * 7)
	now30d := now.Add(-time.Hour * 24 * 30)

	total, err := uc.sg.TotalRedirrectCount(ctx, urlID)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting redirect count: %w", err)
	}

	total24h, err := uc.sg.RedirrectCount(ctx, urlID, now24h, now)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting redirect count: %w", err)
	}

	total7d, err := uc.sg.RedirrectCount(ctx, urlID, now7d, now)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting redirect count: %w", err)
	}

	total30d, err := uc.sg.RedirrectCount(ctx, urlID, now30d, now)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting redirect count: %w", err)
	}

	unique, err := uc.sg.TotalIpCount(ctx, urlID)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting ips count: %w", err)
	}

	unique24h, err := uc.sg.IpCount(ctx, urlID, now24h, now)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting ips count: %w", err)
	}

	unique7d, err := uc.sg.IpCount(ctx, urlID, now7d, now)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting ips count: %w", err)
	}

	unique30d, err := uc.sg.IpCount(ctx, urlID, now30d, now)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting ips count: %w", err)
	}

	return domain.SummaryStatistics{
		Total:     total,
		Total24H:  total24h,
		Total7D:   total7d,
		Total30D:  total30d,
		UniqueIPs: unique,
		Unique24H: unique24h,
		Unique7D:  unique7d,
		Unique30D: unique30d,
	}, nil
}
