package config_test

import (
	"borda/internal/app/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Setenv("APP_ENV", "69")
	os.Setenv("POSTGRES_DB", "crazy_go")
	// Get config
	config := config.Config()

	assert := assert.New(t)
	assert.Equal(8080, config.Get("http.port"))
	assert.Equal("crazy_go", config.Get("db.postgres.name"))
	assert.Equal("69", config.Get("app_env"))
}
