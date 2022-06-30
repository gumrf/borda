package config

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func NewConfig() *Config {
	config := Config{
		Viper: viper.New(),
	}

	viper.AllowEmptyEnv(true)

	// App
	config.setFromEnv("app.env", "BORDA_ENV", "dev")
	config.setFromEnv("app.host", "HOST", "localhost")
	config.setFromEnv("app.port", "PORT", 6900)

	// DB
	config.setFromEnv("db.url", "DATABASE_URL", "")
	config.setFromEnv("db.host", "DB_HOST", "localhost")
	config.setFromEnv("db.port", "DB_PORT", 5432)
	config.setFromEnv("db.user", "DB_USER", "postgres")
	config.setFromEnv("db.password", "DB_PASSWORD", "db_secret_password")
	config.setFromEnv("db.name", "DB_NAME", "borda")
	config.setFromEnv("db.sslmode", "DB_SSLMODE", "disable")
	config.setFromEnv("db.migrations_dir", "DB_MIGRATIONS_DIR", "file://./migrations")

	// jwt
	config.setFromEnv("jwt.signing_key", "JWT_SIGNING_KEY", "jwt_secret_key")
	config.setFromEnv("jwt.key_expire_hours_count", "JWT_EXPIRE_HOURS_COUNT", 24)

	// Other
	config.setFromEnv("password_salt", "PASSWORD_SALT", "pswd_secret_phrase")
	config.setFromEnv("logger_path", "LOGGER_PATH", "logs")

	return &config
}

func (c *Config) setFromEnv(key, env string, defaultValue interface{}) {
	if err := c.BindEnv(key, env); err != nil {
		panic(err)
	}
	c.SetDefault(key, defaultValue)
}

func (c *Config) String() string {
	json, err := json.MarshalIndent(c.AllSettings(), "", "  ")
	if err != nil {
		panic(err)
	}

	// yaml, err := yaml.Marshal(c)
	// if err != nil {
	// 	panic(err)
	// }

	return "Configuration: \n" + string(json)
}

type JWTConfig struct {
	SigningKey string
	ExpireTime time.Duration
}

func (c *Config) JWT() JWTConfig {
	hoursCount := c.GetInt("jwt.key_expire_hours_count")

	return JWTConfig{
		SigningKey: c.GetString("jwt.signing_key"),
		ExpireTime: time.Duration(hoursCount) * time.Hour,
	}
}

func (c *Config) Salt() string {
	return c.GetString("password_salt")
}

// DBURL ...
func (c *Config) DBURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.GetString("db.user"),
		c.GetString("db.password"),
		c.GetString("db.host"),
		c.GetString("db.port"),
		c.GetString("db.name"),
		c.GetString("db.sslmode"),
	)
}

// MigrationsDir ...
func (c *Config) MigrationsDir() string {
	return c.GetString("db.migrations_dir")
}

// LoggerDir ...
func (c *Config) LoggerDir() string {
	return fmt.Sprintf("%s/%s",
		c.GetString("logger.path"),
		c.GetString("logger.file_name"))
}

// TODO: custom config type viper.Viper and expose only get functions.
// TODO: Make config safe for concurrent operations.

// var (
// 	config viper.Viper
// 	once   sync.Once
// )

// func init() {
// 	once.Do(func() {
// 		config = *viper.New()
// 		setDefaults()
// 		loadFromEnv()
// 	})
// }

// func setDefaults() {
// 	config.SetDefault("server.addr", "localhost:8080")

// 	config.SetDefault("db.postgres.host", "localhost")
// 	config.SetDefault("db.postgres.port", 5432)
// 	config.SetDefault("db.postgres.user", "postgres")
// 	config.SetDefault("db.postgres.password", "postgres")
// 	config.SetDefault("db.postgres.name", "postgres")

// 	config.SetDefault("migrations_path", "file://./migrations")

// 	config.SetDefault("logger.path", "logs")
// 	config.SetDefault("logger.file_name", "app.log")

// 	config.SetDefault("jwt.signing_key", "private_key_x69s")
// 	config.SetDefault("jwt.key_expire_hours_count", 24)
// 	config.SetDefault("password_salt", "password_salt_x69s")
// }

// // loadFromEnv reads values from environment variables.
// // If variable isn't bound it will be set to default.
// func loadFromEnv() {
// 	config.AllowEmptyEnv(false)
// 	config.BindEnv("APP_ENV")
// 	config.BindEnv("server.addr", "SERVER_ADDR")

// 	config.BindEnv("db.postgres.host", "POSTGRES_HOST")
// 	config.BindEnv("db.postgres.port", "POSTGRES_PORT")
// 	config.BindEnv("db.postgres.user", "POSTGRES_USER")
// 	config.BindEnv("db.postgres.password", "POSTGRES_PASSWORD")
// 	config.BindEnv("db.postgres.name", "POSTGRES_DB")
// 	//config.BindEnv("db.url", "DATABASE_URL")

// 	config.BindEnv("migrations_path", "MIGRATIONS_PATH")
// }

// func Config() *viper.Viper {
// 	return &config
// }

// func Print() string {
// 	configJSON, err := json.MarshalIndent(config.AllSettings(), "", "  ")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return "Configuration: " + string(configJSON)
// }

// type JWTConfig struct {
// 	SigningKey string
// 	ExpireTime time.Duration
// }

// func JWT() JWTConfig {
// 	expireHoursCount := config.GetInt("jwt.key_expire_hours_count")

// 	return JWTConfig{
// 		SigningKey: config.GetString("jwt.signing_key"),
// 		ExpireTime: time.Duration(expireHoursCount) * time.Hour,
// 	}
// }

// func PasswordSalt() string {
// 	return config.GetString("password_salt")
// }

// func DatabaseURL() string {
// 	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
// 		config.GetString("db.postgres.user"),
// 		config.GetString("db.postgres.password"),
// 		config.GetString("db.postgres.host"),
// 		config.GetInt("db.postgres.port"),
// 		config.GetString("db.postgres.name"))
// }

// // MigrationsPath return path to migrations location
// func MigrationsPath() string {
// 	return config.GetString("migrations_path")
// }

// // LoggerPath return path to migrations location
// func LoggerPath() string {
// 	return fmt.Sprintf("%s/%s",
// 		config.GetString("logger.path"),
// 		config.GetString("logger.file_name"))
// }

// func ServerAddr() string {
// 	return config.GetString("server.addr")
// 	//return fmt.Sprintf("%s:%s",
// 	//	config.GetString("http.host"),
// 	//	config.GetString("http.port"))
// }
