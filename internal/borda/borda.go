package borda

import (
	"borda/internal/api"
	"borda/internal/server"
	"borda/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func Run() {
	// Logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Zap logger initialized")

	// TODO: initialize config

	db, err := postgres.NewPostgresDatabase(os.Getenv("POSTGRES_URI"), "", "")
	if err != nil {
		logger.Error("can't connect to postgres database:", zap.Error(err))
	}

	logger.Info("Connected to Postgres Database")

	fmt.Printf("%+v", db.Stats())

	// TODO: Initialize services

	handlers := api.NewHandlers()

	// HTTP Server
	server := server.NewServer(handlers.Init())

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error occurred while running http server:", zap.Error(err))
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
		logger.Error("failed to stop server", zap.Error(err))
	}

	// Close database connections
	if err := db.Close(); err != nil {
		logger.Error("failed to stop database", zap.Error(err))
	}
}
