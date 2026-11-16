package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urlshortener/internal/cache"
	"urlshortener/internal/config"
	"urlshortener/internal/database"
	"urlshortener/internal/handlers"
	"urlshortener/internal/middlewares"
	"urlshortener/internal/repositories"
	"urlshortener/internal/usecases"

	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		slog.Error("ошибка при работе", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err == nil {
		slog.Info(".env file loaded.")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	// Logger
	logOpts := &slog.HandlerOptions{Level: cfg.LogLevel()}
	logger := slog.New(slog.NewTextHandler(os.Stdout, logOpts))
	slog.SetDefault(logger)

	// DB Connection
	db, err := database.ConnectPostgres(cfg)
	if err != nil {
		return fmt.Errorf("Error on creating Postres connection: %w", err)
	}

	// Migrations
	if err := database.RunMigrations(db); err != nil {
		return fmt.Errorf("Error on running migrations: %w", err)
	}

	// Redis connection
	rdb, err := database.ConnectRedis(cfg)
	if err != nil {
		return fmt.Errorf("Error connecting to Redis: %w", err)
	}

	// Cache and repositories
	redisCache := cache.NewRedisCache(rdb, time.Second)
	urlCache := cache.NewUrlCache(redisCache)
	urlRepository := repositories.NewUrlPostgresRepository(db)
	redirectRepository := repositories.NewRedirectsPostgresRepository(db)
	statisticsRepository := repositories.NewRedirectStatisticsPostgresRepository(db)

	redirectUseCases := usecases.NewRedirectsUseCases(redirectRepository)
	urlUseCases := usecases.NewUrlUseCases(urlRepository, urlCache)
	statisticsUseCases := usecases.NewStatisticsUseCases(statisticsRepository)

	// Handlers
	h := handlers.NewShortenerHandlers(urlUseCases, redirectUseCases, statisticsUseCases)
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux, h)

	// Middlewares
	handler := middlewares.RecoverMiddleware(middlewares.PrometheusMiddleware(mux))
	serverAddr := ":8080"

	server := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	// Starting server
	go func() {
		slog.Info("Server started!", slog.String("Address", serverAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(fmt.Sprintf("Server error: %v", err))
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	slog.Info("Shutting down...")
	if err = server.Shutdown(ctx); err != nil {
		return fmt.Errorf("Error on shutting down: %w", err)
	}

	slog.Info("Server stopped.")
	return nil
}
