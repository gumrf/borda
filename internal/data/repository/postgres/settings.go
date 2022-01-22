package repository

import (
	"borda/internal/core/interfaces"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostgresSettingsRepository struct {
	db                *sqlx.DB
	tableSettingsName string
}

var _ interfaces.SettingsRepository = (*PostgresSettingsRepository)(nil)

func NewPostgresSettingsRepository(db *sqlx.DB) interfaces.SettingsRepository {
	return PostgresSettingsRepository{
		db:                db,
		tableSettingsName: "manage_settings",
	}
}

func (r PostgresSettingsRepository) Get(key string) (value string, err error) {
	query := fmt.Sprintf(`SELECT * FROM public.%s WHERE key=$1`, r.tableSettingsName)

	err = r.db.QueryRowx(query, key).Scan(&value)

	if err == sql.ErrNoRows {
		return value, fmt.Errorf("settings repository get error: Settings not found with key=%v", key)
	}

	if err != nil {
		return value, fmt.Errorf("settings repository get error: %v", err)
	}

	return value, nil
}

func (r PostgresSettingsRepository) Set(key string, value string) (settingsId int, err error) {
	query := fmt.Sprintf(
		`INSERT INTO public.%s (key, value) 
					VALUES ($1, $2)
					ON CONFLICT (key) DO UPDATE 
					SET key = excluded.key, 
					value = excluded.value
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
