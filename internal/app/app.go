package app

import (
	"borda/internal/app/api"
	"borda/internal/app/config"
	"borda/internal/app/logger"
	"fmt"

	pdb "borda/pkg/postgres"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
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

	logger.Log.Info("Connected to DB")

	// Migrations
	if err := Migrate(db); err != nil {
		logger.Log.Fatal("Failed migration: ", err)
	}

	app := fiber.New()
	// Init api routes.
	api.RegisterRoutes(app)

	// Catch OS signals.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		_ = <-quit
		logger.Log.Info("Gracefully shutting down...")
		// Received an interrupt signal, shutdown.
		if err := app.Shutdown(); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Log.Errorf("Oops... Server is not shutting down! Reason: %w", err)
		}
	}()

	// Run server.
	if err := app.Listen(fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)); err != nil {
		logger.Log.Errorf("Oops... Server is not running! Reason: %w", err)
	}

	// Close database connections.
	if err := db.Close(); err != nil {
		logger.Log.Errorf("Oops... Can't close database connections! Reason: %w", err)
	}
}
