package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	_ "rest/docs"
	"rest/internal/cache"
	rd "rest/internal/cache/redis"
	"rest/internal/config"
	"rest/internal/logger"
	"rest/internal/server/handlers/songs/add"
	"rest/internal/server/handlers/songs/get"
	"rest/internal/server/handlers/songs/rm"
	"rest/internal/server/handlers/songs/update"
	"rest/internal/server/handlers/songs/verses"
	"rest/internal/server/middlewares"
	"rest/internal/storage/postgres"
	"syscall"
	"time"
)

func main() {
	const op = "cmd/app.main"
	// init config
	cfg := config.MustLoad()

	// init logger
	log, err := logger.SetupLogger(cfg.Env)
	if errors.Is(err, logger.InvalidEnvErr) {
		fmt.Println(logger.InvalidEnvErr.Error())
		os.Exit(1)
	}

	log.Info("app was started", slog.String("op", op), slog.Any("config", cfg))
	ctx := context.Background()

	//	init storage
	db, err := postgres.New(ctx,
		fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Name),
	)
	if err != nil {
		log.Error("failed to connect to database",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	log.Info("database was initialized", slog.String("op", op))

	// init cache
	rdCache, err := rd.New(ctx, cfg.Cache.Address)
	if err != nil {
		log.Error("failed to connect to cache",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	log.Info("cache was initialized", slog.String("op", op))

	if err = cache.Load(ctx, rdCache, db); err != nil {
		log.Error("failed to load data to cache",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}
	log.Info("data from database was loaded to cache", slog.String("op", op))

	// init router
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(middlewares.New(log))

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Get("/songs", get.New(ctx, log, db))
	router.Get("/verses", verses.New(ctx, log, db, rdCache))
	router.Post("/songs", add.New(ctx, log, db, rdCache))
	router.Put("/songs", update.New(ctx, log, db, rdCache))
	router.Delete("/songs/{id}", rm.New(ctx, log, db, rdCache))

	// start server
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go serverUp(srv, log, signalChan)

	// graceful shutdown
	waitForGracefulShutdown(ctx, srv, db, log, signalChan)
}

func serverUp(srv *http.Server, log *slog.Logger, signalChan chan os.Signal) {
	const op = "cmd/app.ServerUp"
	if err := srv.ListenAndServe(); err != nil {
		log.Error("server error: %s",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		signalChan <- syscall.SIGINT
		return
	}
}

func waitForGracefulShutdown(ctx context.Context, srv *http.Server, db *postgres.Storage, log *slog.Logger,
	signalChan chan os.Signal) {

	const op = "cmd/app.gracefulShutdown"

	<-signalChan
	log.Info("graceful shutdown was started...", slog.String("op", op))
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	db.GetPoolForGracefulShutdown().Close()
	log.Info("all connections in database-pool were closed", slog.String("op", op))
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Error("server shutdown error: %s", slog.String("op", op), slog.String("error", err.Error()))
	}
	log.Info("server shutdown complete")
}
