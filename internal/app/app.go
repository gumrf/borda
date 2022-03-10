package app

import (
	"borda/internal/app/api"
	"borda/internal/app/config"
	"borda/internal/app/logger"
	"borda/internal/data/repository/postgres"

	"fmt"

	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

// Run initializes whole application.
func Run() {
	// Config
	conf := config.Config()

	// Logger
	if err := logger.InitLogger(conf.GetString("logger.path"), conf.GetString("logger.file_name")); err != nil {
		fmt.Println("init logger:", err)
		os.Exit(1)
	}

	// Database
	logger.Log.Info("Connecting to Postgres...: ", config.DatabaseUrl())

	db, err := postgres.Connect(config.DatabaseUrl())
	if err != nil {
		logger.Log.Fatalw("Failed connecting to database:", err)
	}

	app := fiber.New()
	// Init api routes.
	api.RegisterRoutes(app)

	// Catch OS signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Log.Info("Gracefully shutting down...")
		// Received an interrupt signal, shutdown.
		if err := app.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Log.Errorf("Oops... Server is not shutting down! Reason: %w", err)
		}
	}()

	// Run server.
	if err := app.Listen(config.ServerAddr()); err != nil {
		logger.Log.Errorf("Oops... Server is not running! Reason: %w", err)
	}

	// Close database connections.
	if err := db.Close(); err != nil {
		logger.Log.Errorf("Oops... Can't close database connections! Reason: %w", err)
	}
}
