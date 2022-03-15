package app

import (
	"borda/internal/api"
	"borda/internal/config"
	"borda/internal/logger"
	"borda/internal/repository"
	"borda/internal/services"
	"borda/internal/usecase"
	"borda/pkg/hash"
	"borda/pkg/pg"

	"fmt"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

// Run initializes whole application
func Run() {
	conf := config.Config()

	if err := logger.InitLogger(conf.GetString("logger.path"), conf.GetString("logger.file_name")); err != nil {
		fmt.Println("init logger:", err)
		os.Exit(1)
	}

	db, err := pg.Open(config.DatabaseUrl())
	if err != nil {
		logger.Log.Fatalw("Failed to connect to Postgres:", err)
	}
	logger.Log.Info("Connected to Postgres: ", config.DatabaseUrl())

	if err := pg.Migrate(db, config.MigrationsPath()); err != nil {
		logger.Log.Fatalw("Failed to run migrations: %w", err)
	}

	repository := repository.NewRepository(db)
	authService := services.NewAuthService(repository.Users, repository.Teams,
		hash.NewSHA1Hasher(config.PasswordSalt()),
	)

	userUsecase := usecase.NewUserUsecase(repository.Tasks)

	app := fiber.New()

	handlers := api.NewHandler(authService, userUsecase)
	handlers.Init(app)

	// Catch OS signals
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
