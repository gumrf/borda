package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort = "8080"

	Dev  = "dev"
	Prod = "prod"
)

type (
	Config struct {
		Environment string
		Postgres    PostgresConfig
		HTTP        HTTPConfig
	}

	PostgresConfig struct {
		URI      string
		User     string
		Password string
		Database string
	}

	HTTPConfig struct {
		Host string
		Port string
	}
)

func Init(confiDir string) (*Config, error) {
	env := os.Getenv("BORDA_ENV")
	if env == "" {
		env = Dev
	}

	viper.AddConfigPath(confiDir)
	// Additionally check current directory
	viper.AddConfigPath("configs")

	viper.SetConfigType("yaml")
	viper.SetConfigName(env)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w \n", err)
	}

	// Set defaults
	viper.SetDefault("http.port", defaultHTTPPort)

	// Load manually from env
	if err := viper.BindEnv("postgres.uri"); err != nil {
		return nil, fmt.Errorf("viper.BindEnv: %w", err)
	}

	// Fast env
	// viper.AutomaticEnv()

	// Parse config file
	config := new(Config)
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("viper.Unmarshal: %w", err)
	}

	// Todo: fix viper return empty value
	// config.Postgres.URI = viper.GetString("postgres.uri")
	config.Postgres.URI = os.Getenv("POSTGRES_URI")



	return config, nil
}
