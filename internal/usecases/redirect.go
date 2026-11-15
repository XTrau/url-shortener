package usecases

import (
	"context"
	"urlshortener/internal/domain"
)

type RedirectsSaver interface {
	Save(ctx context.Context, redirect domain.RedirectInfo) error
}

type RedirectsUseCases struct {
	redirectsSaver RedirectsSaver
}

func NewRedirectsUseCases(saver RedirectsSaver) *RedirectsUseCases {
	return &RedirectsUseCases{saver}
}

func (uc *RedirectsUseCases) Save(ctx context.Context, redirect domain.RedirectInfo) error {
	return uc.redirectsSaver.Save(ctx, redirect)
}
