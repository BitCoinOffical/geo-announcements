package repo

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
)

type LocationRepo struct {
	db *sql.DB
}

func NewLocationRepo(db *sql.DB) *LocationRepo {
	return &LocationRepo{db: db}
}

func (h *LocationRepo) CreateLocationRepo(ctx context.Context, dto *dto.LocationDTO, userID string) ([]models.DangerousZones, error) {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("tx rollback error: %v", err)
		}
	}()
	row, err := tx.ExecContext(ctx, `
	    INSERT INTO locations (user_id, lat, lon, zone_id)
		SELECT $1, $2, $3, z.zone_id
		FROM zones z
		WHERE ST_Within(ST_SetSRID(ST_MakePoint($3, $2), 4326), z.wkb_geometry)
		LIMIT 1;`, userID, dto.Lat, dto.Lon)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, errors.New("point does not belong to any zone")
	}
	rows, err := tx.QueryContext(ctx, `
	SELECT z.zone_id, z.lat, z.lon,
		ST_Distance(
			z.wkb_geometry,
			ST_SetSRID(ST_Point($1, $2), 4326)
		) AS distance
		FROM zones z
		WHERE z.is_dangerous = TRUE
		ORDER BY z.wkb_geometry <-> ST_SetSRID(ST_Point($1, $2), 4326)
		LIMIT 5;`, dto.Lon, dto.Lat)
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

	return zones, tx.Commit()
}
