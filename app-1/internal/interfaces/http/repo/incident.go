package repo

import (
	"context"
	"database/sql"
	"errors"
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

func (h *IncidentRepo) GetIncidentsRepo(ctx context.Context, limit, offset int) ([]models.Incident, error) {
	query := `SELECT * FROM incidents LIMIT $1 OFFSET $2;`
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
			&incident.Create_at,
			&incident.Update_at,
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
		&incident.Create_at,
		&incident.Update_at,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("incident not found")
		}
		return nil, err
	}
	return &incident, nil

}

func (h *IncidentRepo) GetIncidentStatRepo(ctx context.Context, fromTime *time.Time) (*int, error) {
	query := `SELECT COUNT(DISTINCT user_id) FROM locations WHERE create_at >= $1`
	row := h.db.QueryRowContext(ctx, query, fromTime)
	var count int
	err := row.Scan(
		&count,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("incident not found")
		}
		return nil, err
	}
	return &count, nil
}

func (h *IncidentRepo) CreateIncidentsRepo(ctx context.Context, dto *dto.IncidentDTO) error {
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

func (h *IncidentRepo) UpdateIncidentsByIDRepo(ctx context.Context, dto *dto.IncidentDTO, id int) error {
	query := `UPDATE incidents SET lat = $1, lon = $2, update_at = NOW() WHERE incident_id = $3`
	_, err := h.db.ExecContext(ctx, query, dto.Lat, dto.Lon, id)
	return err
}

func (h *IncidentRepo) DeleteIncidentsByIDRepo(ctx context.Context, id int) error {
	query := `DELETE FROM incidents WHERE incident_id = $1`
	_, err := h.db.ExecContext(ctx, query, id)
	return err
}
