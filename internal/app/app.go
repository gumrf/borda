package app

import (
	"borda/internal/app/api"
	"borda/internal/app/config"
	"borda/internal/app/logger"
	"borda/internal/app/server"
	"fmt"

	pdb "borda/pkg/postgres"
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

// Run initializes whole application.
func Run() {
	// Config
	cfg := config.InitConfig()
	fmt.Println("init config: OK")

	// Logger
	if err := logger.InitLogger(cfg.Additional.LogDir, cfg.Additional.LogFileName); err != nil {
		fmt.Println("init logger:", err)
		os.Exit(1)
	}
	fmt.Println("init logger: OK")
	logger.Log.Info("Logs path: ", filepath.Join(cfg.Additional.LogDir, cfg.Additional.LogFileName))

	// Database
	logger.Log.Info("Database URI: ", cfg.DatabaseURI())

	db, err := pdb.NewConnection(cfg.DatabaseURI())
	if err != nil {
		logger.Log.Fatalw("Failed connecting to database:", err)
	}

	if err := Migrate(db, cfg.DatabaseURI(), cfg.Additional.MigrationsDirName); err != nil {
		logger.Log.Fatal("Failed migration: ", err)
	}

	// Api handlers
	handler := api.NewRoutes()

	// HTTP Server
	server := server.NewServer(handler, cfg.HTTP)

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatal("Error occurred while running http server:", err)
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
		logger.Log.Error("Failed to stop server", err)
	}

	// Close database connections
	if err := db.Close(); err != nil {
		logger.Log.Error("Failed to stop database", err)
	}
}
