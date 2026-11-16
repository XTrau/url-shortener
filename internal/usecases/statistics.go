package usecases

import (
	"context"
	"fmt"
	"time"
	"urlshortener/internal/domain"
)

type StatsGetter interface {
	// Возвращает количество редиректов после t.
	//
	// Если t=nil то возвращает количество за все время.
	RedirrectCount(ctx context.Context, urlID int, t *time.Time) (int, error)
	// Возвращает количество уникальных ip после t.
	//
	// Если t=nil то возвращает количество за все время.
	IpCount(ctx context.Context, urlID int, t *time.Time) (int, error)
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

	total, err := uc.sg.RedirrectCount(ctx, urlID, nil)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting redirect count: %w", err)
	}

	total24h, err := uc.sg.RedirrectCount(ctx, urlID, &now24h)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting redirect count: %w", err)
	}

	total7d, err := uc.sg.RedirrectCount(ctx, urlID, &now7d)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting redirect count: %w", err)
	}

	total30d, err := uc.sg.RedirrectCount(ctx, urlID, &now30d)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting redirect count: %w", err)
	}

	unique, err := uc.sg.IpCount(ctx, urlID, nil)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting ips count: %w", err)
	}

	unique24h, err := uc.sg.IpCount(ctx, urlID, &now24h)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting ips count: %w", err)
	}

	unique7d, err := uc.sg.IpCount(ctx, urlID, &now7d)
	if err != nil {
		return domain.SummaryStatistics{}, fmt.Errorf("Error on getting ips count: %w", err)
	}

	unique30d, err := uc.sg.IpCount(ctx, urlID, &now30d)
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
