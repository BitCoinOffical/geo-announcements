package postgres

import (
	"database/sql"
	"fmt"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg *config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME))
	if err != nil {
		return nil, err
	}
	return db, nil
}
