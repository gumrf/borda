package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// TODO: custom config type viper.Viper and expose only get functions.
// TODO: Make config safe for concurrent operations.

var (
	// config is config
	config viper.Viper
	once   sync.Once
)

func init() {
	once.Do(func() {
		config = *viper.New()
		setDefaults()
		loadFromEnv()
	})
}

// setDefaults sets default values :)
func setDefaults() {
	config.SetDefault("http.host", "localhost")
	config.SetDefault("http.port", 8080)
	
	config.SetDefault("db.postgres.host", "localhost")
	config.SetDefault("db.postgres.port", 5432)
	config.SetDefault("db.postgres.user", "postgres")
	config.SetDefault("db.postgres.password", "postgres")
	config.SetDefault("db.postgres.name", "postgres")
	config.SetDefault("db.postgres.migrations_path", "file://../migrations/")

	config.SetDefault("logger.path", "logs")
	config.SetDefault("logger.file_name", "app.log")

	config.SetDefault("jwt.signing_key", "private_key_x69s")
	config.SetDefault("jwt.key_expire_hours_count", 24)
	config.SetDefault("password_salt", "password_salt_x69s")
}

// loadFromEnv reads values from environment variables.
// If variable isn't bound it will be set to default.
func loadFromEnv() {
	config.AllowEmptyEnv(false)
	config.BindEnv("APP_ENV")
	config.BindEnv("db.postgres.host", "POSTGRES_HOST")
	config.BindEnv("db.postgres.port", "POSTGRES_PORT")
	config.BindEnv("db.postgres.user", "POSTGRES_USER")
	config.BindEnv("db.postgres.password", "POSTGRES_PASSWORD")
	config.BindEnv("db.postgres.name", "POSTGRES_DB")
	config.BindEnv("db.postgres.migrations_path", "POSTGRES_MIGRATION_PATH")
}

func Config() *viper.Viper {
	return &config
}

func Print() string {
	configJSON, err := json.MarshalIndent(config.AllSettings(), "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return "Configuration: " + string(configJSON)
}

type JWTConfig struct {
	SigningKey string
	ExpireTime time.Duration
}

func JWT() JWTConfig {
	expireHoursCount := config.GetInt("jwt.key_expire_hours_count")

	return JWTConfig{
		SigningKey: config.GetString("jwt.signing_key"),
		ExpireTime: time.Duration(expireHoursCount) * time.Hour,
	}
}

func PasswordSalt() string {
	return config.GetString("auth.password_salt")
}

// DatabaseUrl returns full Postgres connection url.
func DatabaseUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.GetString("db.postgres.user"),
		config.GetString("db.postgres.password"),
		config.GetString("db.postgres.host"),
		config.GetInt("db.postgres.port"),
		config.GetString("db.postgres.name"))
}

// MigrationsPath return path to migrations location
func MigrationsPath() string {
	return config.GetString("db.postgres.migrations_path")
}

// LoggerPath return path to migrations location
func LoggerPath() string {
	return fmt.Sprintf("%s/%s",
		config.GetString("logger.path"),
		config.GetString("logger.file_name"))
}

func ServerAddr() string {
	return fmt.Sprintf("%s:%s",
		config.GetString("http.host"),
		config.GetString("http.port"))
}
