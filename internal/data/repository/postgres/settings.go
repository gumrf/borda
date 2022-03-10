package postgres

import (
	"borda/internal/core"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SettingsRepository struct {
	db                *sqlx.DB
	tableSettingsName string
}

var _ core.SettingsRepository = (*SettingsRepository)(nil)

func NewSettingsRepository(db *sqlx.DB) core.SettingsRepository {
	return SettingsRepository{
		db:                db,
		tableSettingsName: "settings",
	}
}

func (r SettingsRepository) Get(key string) (value string, err error) {
	query := fmt.Sprintf(`
		SELECT *
		FROM public.%s
		WHERE key=$1
		LIMIT 1`, r.tableSettingsName)

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

func (r SettingsRepository) Set(key string, value string) (settingsId int, err error) {
	query := fmt.Sprintf(`
		INSERT INTO public.%s (
			key,
			value
		) 
		VALUES ($1, $2)
		ON CONFLICT (key) DO UPDATE 
		SET key = excluded.key, value = excluded.value
		RETURNING id`,
		r.tableSettingsName)
	id := -1
	err = r.db.QueryRowx(query, value, key).Scan(&id)

	if err != nil || id == -1 {
		return id, fmt.Errorf("team repository create error: %v", err)
	}

	if err != nil {
		return id, fmt.Errorf("settings repository get error: %v", err)
	}

	return id, nil
}
