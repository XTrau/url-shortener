package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"
	"urlshortener/internal/config"
	"urlshortener/internal/database"
	"urlshortener/internal/domain"
	"urlshortener/internal/repositories"
	"urlshortener/internal/usecases"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("error on loading config", slog.String("err", err.Error()))
		os.Exit(1)
	}

	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		slog.Error("error on connecting db", slog.String("err", err.Error()))
		os.Exit(1)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	repo := repositories.NewRedirectsPostgresRepository(db)
	redirectUseCases := usecases.NewRedirectsUseCases(repo)

	now := time.Now()
	end := now.Add(-time.Hour * 24 * 60)
	for now.After(end) {
		req := (70 + rand.Intn(30))
		for range req {
			addr := fmt.Sprintf("%v.%v.%v.%v", 1+rand.Intn(5), 1+rand.Intn(6), 1+rand.Intn(5), 1+rand.Intn(5))

			redirect := domain.RedirectInfo{
				CreatedAt: now,
				IP:        addr,
				UserAgent: "",
				UrlID:     3,
			}
			redirectUseCases.Save(context.Background(), redirect)
		}
		fmt.Println("Added", now)

		now = now.Add(-time.Hour)
	}
}
