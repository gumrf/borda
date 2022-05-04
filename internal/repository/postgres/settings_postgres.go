package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SettingsRepository struct {
	db *sqlx.DB
}

func NewSettingsRepository(db *sqlx.DB) *SettingsRepository {
	return &SettingsRepository{db: db}
}

func (r SettingsRepository) Get(key string) (value string, err error) {
	query := fmt.Sprintf(`
		SELECT value
		FROM public.%s
		WHERE key=$1
		LIMIT 1`, settingsTable)

	err = r.db.QueryRowx(query, key).Scan(&value)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return value, fmt.Errorf("settings repository get error: Settings not found with key=%v", key)
		default:
			return value, fmt.Errorf("settings repository get error: %v", err)
		}
	}

	return value, nil
}

func (r SettingsRepository) Set(key string, value string) (settingId int, err error) {
	query := fmt.Sprintf(`
		INSERT INTO public.%s (
			key,
			value
		)
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE
		SET key = excluded.key, value = excluded.value
		RETURNING id`,
		settingsTable)
	id := -1
	err = r.db.QueryRowx(query, key, value).Scan(&id)

	if err != nil || id == -1 {
		return id, fmt.Errorf("team repository create error: %v", err)
	}

	if err != nil {
		return id, fmt.Errorf("settings repository get error: %v", err)
	}

	return id, nil
}
