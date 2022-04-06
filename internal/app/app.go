package app

import (
	"borda/internal/api"
	"borda/internal/config"
	"borda/internal/logger"
	"borda/internal/repository"
	"borda/internal/service"
	"borda/pkg/hash"
	"borda/pkg/pg"

	"fmt"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
)

// @title                       CTF Borda API
// @version                     0.1
// @description                 REST API for CTF Borda.
// @host                        localhost:8080
// @BasePath                    /api/v1
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

// Run initializes whole application
func Run() {
	conf := config.Config()

	if err := logger.InitLogger(conf.GetString("logger.path"), conf.GetString("logger.file_name")); err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}

	db, err := pg.Open(config.DatabaseURL())
	if err != nil {
		logger.Log.Fatalw("Failed to connect to Postgres:", err)
	}
	logger.Log.Info("Connected to Postgres: ", config.DatabaseURL())

	if err := pg.Migrate(db, config.MigrationsPath()); err != nil {
		logger.Log.Fatalf("Failed to run migrations: %v", err)
	}

	// Repository
	repository := repository.NewRepository(db)

	// Services
	authService := service.NewAuthService(repository.Users, repository.Teams,
		hash.NewSHA1Hasher(config.PasswordSalt()),
	)
	userService := service.NewUserService(repository.Users, repository.Tasks, repository.Teams)
	adminService := service.NewAdminService(repository.Tasks)

	//Usecases
	// usecaseGetAllUsers := usecase.NewUsecaseGetAllUsers(repository.Users, repository.Teams)

	app := fiber.New()

	// Handlers
	handlers := api.NewHandler(authService, userService, adminService, repository)
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
