package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (h *IncidentRepo) GetTop(ctx context.Context, limit int) ([]models.Incident, error) {
	query := `SELECT * FROM incidents WHERE status = 'public' ORDER BY incident_id ASC LIMIT $1`
	rows, err := h.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("db.QueryContext: %w", err)
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
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		incidents = append(incidents, incident)
	}
	return incidents, nil
}

func (h *IncidentRepo) GetIncidents(ctx context.Context, limit, offset int) ([]models.Incident, error) {
	query := `SELECT * FROM incidents WHERE status = 'public' ORDER BY incident_id ASC LIMIT $1 OFFSET $2;`
	rows, err := h.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("db.QueryContext: %w", err)
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
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		incidents = append(incidents, incident)
	}
	return incidents, nil
}

func (h *IncidentRepo) GetIncidentByID(ctx context.Context, id int) (*models.Incident, error) {
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
		return nil, fmt.Errorf("rows.Scan: %w", err)
	}
	return &incident, nil

}

func (h *IncidentRepo) GetIncidentStat(ctx context.Context, fromTime *time.Time) (*models.UsersInDangerousZones, error) {
	query := `SELECT COUNT(DISTINCT user_id), zone_id FROM user_checks_loс WHERE create_at >= $1 GROUP BY zone_id`
	row := h.db.QueryRowContext(ctx, query, fromTime)
	var users models.UsersInDangerousZones
	err := row.Scan(
		&users.UserCount,
		&users.Zone_id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("incident not found: ")
		}
		return nil, fmt.Errorf("rows.Scan: %w", err)
	}
	return &users, nil
}

func (h *IncidentRepo) CreateIncidents(ctx context.Context, dto *dto.IncidentDTO) error {
	query := `INSERT INTO incidents (lat, lon) VALUES ($1, $2)`
	res, err := h.db.ExecContext(ctx, query, dto.Lat, dto.Lon)
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return errors.New("failed create incident")
	}
	return nil
}

func (h *IncidentRepo) UpdateZones(ctx context.Context, dto *dto.IncidentDTO) error {
	queue := `UPDATE zones SET is_dangerous = TRUE WHERE zone_id = (SELECT z.zone_id FROM zones z WHERE ST_Contains(
			z.wkb_geometry,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)
			) LIMIT 1;`
	_, err := h.db.ExecContext(ctx, queue, dto.Lat, dto.Lon)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}
	return nil
}

func (h *IncidentRepo) UpdateIncidentsByID(ctx context.Context, dto *dto.IncidentDTO, id int) error {
	query := `UPDATE incidents SET lat = $1, lon = $2, update_at = NOW() WHERE incident_id = $3`
	_, err := h.db.ExecContext(ctx, query, dto.Lat, dto.Lon, id)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}
	return nil
}

func (h *IncidentRepo) DeleteIncidentsByID(ctx context.Context, id int) error {
	query := `UPDATE incidents SET status = 'hide', deleted_at = NOW() WHERE incident_id = $1`
	_, err := h.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}
	return nil
}
