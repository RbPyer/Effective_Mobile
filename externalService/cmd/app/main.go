package main

import (
	"context"
	_ "externalService/docs"
	"externalService/internal/config"
	"externalService/internal/logger"
	"externalService/internal/server/handlers/info"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	const op = "cmd/app.main"

	cfg := config.MustLoad()

	log, err := logger.SetupLogger(cfg.Env)
	if err != nil {
		fmt.Printf("failed to setup logger: %v\n", err)
		os.Exit(1)
	}

	log.Info("starting external service")

	var counter atomic.Uint32

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Get("/info", info.New(&counter, log))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		if err = srv.ListenAndServe(); err != nil {
			log.Error("server listening error", slog.String("op", op), slog.String("error", err.Error()))
			return
		}
	}()
	log.Info("server is running", slog.String("op", op))

	<-signalChan
	log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Error("server shutdown error", slog.String("op", op), slog.String("error", err.Error()))
	}
}
