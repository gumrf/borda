package borda

import (
	"borda/internal/app/api"
	"borda/internal/app/server"
	"borda/internal/app/setup"
	"borda/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func Run() {

	// Config & logger
	cfg := setup.InitConfig()
	fmt.Println("[+] CONFIG: OK")

	logger, err := setup.InitLogger(cfg.Additional.LogDir, cfg.Additional.LogFileName)
	if err != nil {
		fmt.Println("[-] Error on init logger:", err)
		os.Exit(1)
	}
	logger.Info("[+] LOGS: ", filepath.Join(cfg.Additional.LogDir, cfg.Additional.LogFileName))
	
	// Database
	logger.Info("[+] DATABASE URI: ", cfg.DatabaseURI())
	
	db, err := postgres.NewPostgresDatabase(cfg.DatabaseURI())
	if err != nil {
		logger.Fatalw("Failed connecting to database:", err)
	}
	logger.Info("[+] CONNECT TO DB: OK")

	// // TODO: Initialize services

	// Api handlers
	handler := api.NewRoutes()

	// HTTP Server
	server := server.NewServer(handler, cfg.HTTP)

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Error occurred while running http server:", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// wait for signal
	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := server.Stop(ctx); err != nil {
		logger.Fatal("Failed to stop server", err)
	}

	// Close database connections
	if err := db.Close(); err != nil {
		logger.Fatal("Failed to stop database", err)
	}

	os.Exit(1)
}
