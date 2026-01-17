package postgres

import (
	"database/sql"
	"fmt"

	"github.com/BitCoinOffical/geo-announcements/app-1/config"
	_ "github.com/lib/pq"
)

func NewPostgres(cfg *config.PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBport, cfg.DBUser, cfg.DBPassword, cfg.DBName))
	if err != nil {
		return nil, err
	}
	return db, nil
}
