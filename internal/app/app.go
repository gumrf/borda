package app

import (
	"borda/internal/app/api"
	"borda/internal/app/server"
	"borda/internal/app/setup"

	pdb "borda/pkg/postgres"
	"context"
	"errors"
	"log"
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
	log.Println("Config initialized")

	if err := setup.InitLogger(cfg.Additional.LogDir, cfg.Additional.LogFileName); err != nil {
		log.Println("error on init logger:", err)
		os.Exit(1)
	}

	logger := setup.GetLogger()

	logger.Info("Logs path: ", filepath.Join(cfg.Additional.LogDir, cfg.Additional.LogFileName))

	// Database
	logger.Info("Database URI: ", cfg.DatabaseURI())

	db, err := pdb.NewConnection(cfg.DatabaseURI())
	if err != nil {
		logger.Fatalw("Failed connecting to database:", err)
	}
	logger.Info("Connected to DB")

	if err := Migrate(db, cfg.DatabaseURI(), cfg.Additional.MigrationsDirName); err != nil {
		logger.Fatal("Failed migration: ", err)
	}
	logger.Info("Migration did run successfully")

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
}
