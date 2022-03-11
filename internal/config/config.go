package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

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
	config.SetDefault("auth.jwt.signing_key", "sd14fs88ef123dsD0101KdlfICDpdsdfsd435csd")
	config.SetDefault("auth.jwt.expire_time_hours", 24)
	config.SetDefault("auth.password_salt", "random_string")
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

//GetJwtEntity returns signing key and expire timr
func GetJwtEntity() (string, int) {
	return config.GetString("auth.jwt.signing_key"), config.GetInt("auth.jwt.expire_time_hours")
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
