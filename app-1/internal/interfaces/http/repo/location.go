package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
)

type LocationRepo struct {
	db *sql.DB
}

func NewLocationRepo(db *sql.DB) *LocationRepo {
	return &LocationRepo{db: db}
}

func (h *LocationRepo) CreateLocationRepo(ctx context.Context, dto *dto.LocationDTO, userID string) error {
	query := `INSERT INTO locations (user_id, lat, lon ) VALUES ($1, $2, $3)`
	res, err := h.db.ExecContext(ctx, query, userID, dto.Lat, dto.Lon)
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return errors.New("failed create location")
	}
	return nil
}
