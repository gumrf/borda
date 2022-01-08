package setup

import (
	"fmt"
	"os"
)


type (
	Config struct {
		Postgres    PostgresConfig
		HTTP        HTTPConfig
		Additional  AdditionalConfig
		Credentials CredentialsConfig
 	}

	CredentialsConfig struct {
		JWTKey string
	}

	AdditionalConfig struct {
		LogDir      string
		LogFileName string
		JWTExpTime  int
	}

	PostgresConfig struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
	}

	HTTPConfig struct {
		Host string
		Port string
	}
)

func (c *Config) DatabaseURI() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.Postgres.User, c.Postgres.Password, c.Postgres.Host, c.Postgres.Port, c.Postgres.Database)
}

func getEnv(name string, missing string) string {
	env := os.Getenv(name)
	if env == "" {
		return missing
	}
	return env
}

func InitConfig() *Config {
	config := new(Config)
	config.Postgres.Host = getEnv("DB_HOST", "localhost")
	config.Postgres.Port = getEnv("DB_PORT", "5432")
	config.Postgres.User = getEnv("DB_USER", "postgres")
	config.Postgres.Password = getEnv("DB_PASSWORD", "postgres")
	config.Postgres.Database = getEnv("DB_NAME", "postgres")

	config.HTTP.Host = getEnv("HTTP_HOST", "localhost")
	config.HTTP.Port = getEnv("HTTP_PORT", "8080")

	config.Additional.JWTExpTime = 1800
	config.Additional.LogDir = "logs"
	config.Additional.LogFileName = "logs.txt"

	config.Credentials.JWTKey = getEnv("JWT_KEY", "jwt_jwt_key_change_meee")

	return config
}
