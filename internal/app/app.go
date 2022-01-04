package borda

import (
	"borda/internal/app/api"
	"borda/internal/app/config"
	"borda/internal/app/server"
	"borda/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func Run(configPath string) {
	// Logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("Failed to initialize zap logger: %v", err))
	}

	logger.Info("Logger initialized")

	// Config
	cfg, err := config.Init(configPath)
	if err != nil {
		logger.Error("Failed to initialize config", zap.Error(err))
	}

	fmt.Println(cfg.Postgres.URI)

	// Database
	logger.Info("Connecting to database")

	db, err := postgres.NewPostgresDatabase(cfg.Postgres.URI, "", "")
	if err != nil {
		logger.Error("Failed connecting to database:", zap.Error(err))
	}

	logger.Info("Connected to database", zap.String("uri", db.DataSourceName))

	// TODO: Initialize services

	// Api handlers
	handler := api.NewApiHandler(logger)

	// HTTP Server
	server := server.NewServer(handler.Init(), logger)

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Error occurred while running http server:", zap.Error(err))
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// wait for signal
	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := server.Stop(ctx); err != nil {
		logger.Error("Failed to stop server", zap.Error(err))
	}

	// Close database connections
	if err := db.Close(); err != nil {
		logger.Error("Failed to stop database", zap.Error(err))
	}

	logger.Info("Server stoped")

	if err := logger.Sync(); err != nil {
		panic(fmt.Errorf("logger.Sync: %w", err))
	}

	os.Exit(1)
}
