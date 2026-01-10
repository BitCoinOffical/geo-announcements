package repo

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
)

type IncidentRepo struct {
	db *sql.DB
}

func NewIncidentRepo(db *sql.DB) *IncidentRepo {
	return &IncidentRepo{db: db}
}

func (h *IncidentRepo) GetTopRepo(ctx context.Context, limit int) ([]models.Incident, error) {
	query := `SELECT * FROM incidents WHERE status = 'public' ORDER BY incident_id ASC LIMIT $1`
	rows, err := h.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	var incidents []models.Incident
	for rows.Next() {
		var incident models.Incident
		if err := rows.Scan(
			&incident.Incident_id,
			&incident.Lat,
			&incident.Lon,
			&incident.Status,
			&incident.Create_at,
			&incident.Update_at,
			&incident.Deleted_at,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}
	return incidents, nil
}

func (h *IncidentRepo) GetIncidentsRepo(ctx context.Context, limit, offset int) ([]models.Incident, error) {
	query := `SELECT * FROM incidents WHERE status = 'public' ORDER BY incident_id ASC LIMIT $1 OFFSET $2;`
	rows, err := h.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	var incidents []models.Incident
	for rows.Next() {
		var incident models.Incident
		if err := rows.Scan(
			&incident.Incident_id,
			&incident.Lat,
			&incident.Lon,
			&incident.Status,
			&incident.Create_at,
			&incident.Update_at,
			&incident.Deleted_at,
		); err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}
	return incidents, nil
}

func (h *IncidentRepo) GetIncidentByIDRepo(ctx context.Context, id int) (*models.Incident, error) {
	query := `SELECT * FROM incidents WHERE incident_id = $1`
	row := h.db.QueryRowContext(ctx, query, id)
	var incident models.Incident
	err := row.Scan(
		&incident.Incident_id,
		&incident.Lat,
		&incident.Lon,
		&incident.Status,
		&incident.Create_at,
		&incident.Update_at,
		&incident.Deleted_at,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("incident not found")
		}
		return nil, err
	}
	return &incident, nil

}

func (h *IncidentRepo) GetIncidentStatRepo(ctx context.Context, fromTime *time.Time) (*models.UsersInDangerousZones, error) {
	query := `SELECT COUNT(DISTINCT user_id), zone_id FROM user_checks_loс WHERE create_at >= $1 GROUP BY zone_id`
	row := h.db.QueryRowContext(ctx, query, fromTime)
	var users models.UsersInDangerousZones
	err := row.Scan(
		&users.UserCount,
		&users.Zone_id,
	)
	if err != nil {
		log.Println(err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("incident not found: ")
		}
		return nil, err
	}
	return &users, nil
}

func (h *IncidentRepo) CreateIncidentsRepo(ctx context.Context, dto *dto.IncidentDTO) error {
	tx, err := h.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("tx rollback error: %v", err)
		}
	}()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO incidents (lat, lon) VALUES ($1, $2)`,
		dto.Lat,
		dto.Lon,
	)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(
		ctx,
		`
		UPDATE zones
		SET is_dangerous = TRUE
		WHERE zone_id = (
			SELECT z.zone_id
			FROM zones z
			WHERE ST_Contains(
				z.wkb_geometry,
				ST_SetSRID(ST_Point($1, $2), 4326)
			)
			LIMIT 1
		)
		`,
		dto.Lon,
		dto.Lat,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (h *IncidentRepo) UpdateIncidentsByIDRepo(ctx context.Context, dto *dto.IncidentDTO, id int) error {
	query := `UPDATE incidents SET lat = $1, lon = $2, update_at = NOW() WHERE incident_id = $3`
	_, err := h.db.ExecContext(ctx, query, dto.Lat, dto.Lon, id)
	return err
}

func (h *IncidentRepo) DeleteIncidentsByIDRepo(ctx context.Context, id int) error {
	query := `UPDATE incidents SET status = 'hide', deleted_at = NOW() WHERE incident_id = $1`
	_, err := h.db.ExecContext(ctx, query, id)
	return err
}
