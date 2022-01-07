package borda

import (
	"borda/internal/app/api"
	"borda/internal/app/server"
	"borda/internal/app/setup"
	"borda/pkg/postgres"
	"context"
	"errors"
	"os"
	"path/filepath"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	// Config & logger
	cfg := setup.InitConfig()
	fmt.Println("[+] CONFIG: OK")
	err := setup.InitLogger(cfg.Additional.LogDir, cfg.Additional.LogFileName)
	if err != nil {
		fmt.Println("[-] Error on init logger:", err)
		os.Exit(1)
	}
	fmt.Println("[+] LOGS: ", filepath.Join(cfg.Additional.LogDir, cfg.Additional.LogFileName) )

	logger, err := setup.GetReadyLogger()
	if err != nil {
		fmt.Println("[-] Error on get logger:", err)
	}

	// logger.Println("Basic app init!")

	fmt.Println("[+] DATABASE URI: ", cfg.DB_URI())

	// Database
	db, err := postgres.NewPostgresDatabase(cfg.DB_URI(), "", "")
	if err != nil {
		logger.Fatalln("Failed connecting to database:", err)
	}
	fmt.Println("[+] CONNECT TO DB: OK")

	// // TODO: Initialize services

	// Api handlers
	handler := api.NewRoutes()

	// HTTP Server
	server := server.NewServer(handler, cfg.HTTP.Host, cfg.HTTP.Port)

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalln("Error occurred while running http server:", err)
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
		logger.Fatalln("Failed to stop server", err)
	}

	// Close database connections
	if err := db.Close(); err != nil {
		logger.Fatalln("Failed to stop database", err)
	}

	os.Exit(1)
}
