package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
)

type LocationRepo struct {
	db *sql.DB
}

func NewLocationRepo(db *sql.DB) *LocationRepo {
	return &LocationRepo{db: db}
}

func (h *LocationRepo) CreateLocation(ctx context.Context, dto *dto.LocationDTO, userID string) error {
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

func (h *LocationRepo) GetDangerZones(ctx context.Context, dto *dto.LocationDTO, userID string) ([]models.DangerousZones, error) {
	query := `
	SELECT z.zone_id, z.lat, z.lon,
		ST_Distance(
			z.wkb_geometry,
			ST_SetSRID(ST_Point($1, $2), 4326)
		) AS distanc FROM zones z
		WHERE z.is_dangerous = TRUE
		ORDER BY z.wkb_geometry <-> ST_SetSRID(ST_Point($1, $2), 4326) LIMIT 5;`
	rows, err := h.db.QueryContext(ctx, query, userID, dto.Lat, dto.Lon)
	if err != nil {
		return nil, err
	}
	var zones []models.DangerousZones
	for rows.Next() {
		var zone models.DangerousZones
		if err := rows.Scan(
			&zone.Zone_id,
			&zone.Lat,
			&zone.Lon,
			&zone.Distant,
		); err != nil {
			return nil, err
		}
		zones = append(zones, zone)
	}

	return zones, nil
}
