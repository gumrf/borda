package postgres_test

import (
	"borda/internal/repository"
	"borda/internal/repository/postgres"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewSettingsRepository(t *testing.T) {
	db := MustOpenDB(t)

	settingsRepository := postgres.NewSettingsRepository(db)

	require.Implements(t, (*repository.SettingsRepository)(nil), settingsRepository)
}

func Test_SettingsRepository_Set(t *testing.T) {
	db := MustOpenDB(t)
	repo := postgres.NewSettingsRepository(db)
	assert := assert.New(t)

	key := "team_limit"
	value := "5"
	settingsId := 1

	id, err := repo.Set(key, value)
	if err != nil {
		t.Fatalf("Set error: %v\n", err)
	}

	assert.Equal(id, settingsId, "they should be equal")
	assert.NotNil(id, "must be not nil")
}

func Test_SettingsRepository_Get(t *testing.T) {
	db := MustOpenDB(t)
	repo := postgres.NewSettingsRepository(db)
	assert := assert.New(t)

	key := "team_limit"
	value := "5"

	_, err := repo.Set(key, value)
	if err != nil {
		t.Fatalf("Test settings asserted not created: %v\n", err)

	}

	testValue, err := repo.Get(key)
	// if err != nil {
	// 	t.Fatalf("Test settings asserted not created: %v\n", err)
	// }

	assert.Error(err, "error should be nil")
	assert.Equal(value, testValue, fmt.Sprintf("value should be equal <%v>, not <%v>", value, testValue))

	_, err = repo.Get("1337")
	assert.Error(err, "Settings not found with key=1337", "they should be equal")

}
