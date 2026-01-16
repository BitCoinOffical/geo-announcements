package repo_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/models"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/repo"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	limit                  = 100
	page                   = 1
	id                     = 0
	StatsTimeWindowMinutes = 5
)

func TestGetTop(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repo.NewIncidentRepo(db)
	q := `SELECT * FROM incidents WHERE status = 'public' ORDER BY incident_id ASC LIMIT $1`
	tests := []struct {
		name  string
		query string
		err   error
		ok    bool
	}{
		{name: "ok", query: q, err: nil, ok: true},
		{name: "error", query: q, err: errors.New("db down"), ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(tt.query)).WithArgs(limit).WillReturnError(tt.err).WillReturnRows(sqlmock.NewRows([]string{
				"incident_id",
				"lat",
				"lon",
				"status",
				"create_at",
				"update_at",
				"deleted_at"}))
			rows, err := repo.GetTop(t.Context(), limit)
			if tt.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

			assert.Len(t, rows, 0)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}

}

func TestGetIncidents(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repo.NewIncidentRepo(db)
	q := `SELECT * FROM incidents WHERE status = 'public' ORDER BY incident_id ASC LIMIT $1 OFFSET $2`
	tests := []struct {
		name  string
		query string
		err   error
		ok    bool
	}{
		{name: "ok", query: q, err: nil, ok: true},
		{name: "error", query: q, err: errors.New("db down"), ok: false},
	}
	offset := (page - 1) * limit
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(tt.query)).WithArgs(limit, offset).WillReturnError(tt.err).WillReturnRows(sqlmock.NewRows([]string{
				"incident_id",
				"lat",
				"lon",
				"status",
				"create_at",
				"update_at",
				"deleted_at"}))
			rows, err := repo.GetIncidents(t.Context(), limit, offset)
			if tt.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

			assert.Len(t, rows, 0)
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}

func TestGetIncidentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repo.NewIncidentRepo(db)
	q := `SELECT * FROM incidents WHERE incident_id = $1`
	tests := []struct {
		name  string
		query string
		err   error
		ok    bool
	}{
		{name: "ok", query: q, err: nil, ok: true},
		{name: "error", query: q, err: errors.New("db down"), ok: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(tt.query)).WithArgs(id).WillReturnError(tt.err).WillReturnRows(sqlmock.NewRows([]string{
				"incident_id",
				"lat",
				"lon",
				"status",
				"create_at",
				"update_at",
				"deleted_at"}).AddRow(
				id, 10.0, 20.0, "open", nil, nil, nil,
			))
			rows, err := repo.GetIncidentByID(t.Context(), id)
			if tt.ok {
				assert.NoError(t, err)
				assert.Equal(t, &models.Incident{Incident_id: 0, Lat: 10, Lon: 20, Status: "open", Create_at: nil, Update_at: nil, Deleted_at: nil}, rows)
			} else {
				assert.Error(t, err)
				assert.Nil(t, rows)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}

func TestGetIncidentStat(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repo.NewIncidentRepo(db)
	q := `SELECT COUNT(DISTINCT user_id), zone_id FROM user_checks_loс WHERE create_at >= $1 GROUP BY zone_id`
	tests := []struct {
		name  string
		query string
		err   error
		ok    bool
	}{
		{name: "ok", query: q, err: nil, ok: true},
		{name: "error", query: q, err: errors.New("db down"), ok: false},
	}
	fromTime := time.Now().Add(-time.Duration(StatsTimeWindowMinutes) * time.Minute)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(tt.query)).WithArgs(fromTime).WillReturnError(tt.err).WillReturnRows(sqlmock.NewRows([]string{
				"Zone_id",
				"UserCount"}).AddRow(1, 1))
			rows, err := repo.GetIncidentStat(t.Context(), &fromTime)
			if tt.ok {
				assert.NoError(t, err)
				assert.Equal(t, &models.UsersInDangerousZones{Zone_id: 1, UserCount: 1}, rows)
			} else {
				assert.Error(t, err)
				assert.Nil(t, rows)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}

func TestCreateIncidents(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repo.NewIncidentRepo(db)
	q := `INSERT INTO incidents (lat, lon) VALUES ($1, $2)`
	tests := []struct {
		name  string
		query string
		err   error
		ok    bool
	}{
		{name: "ok", query: q, err: nil, ok: true},
		{name: "error", query: q, err: errors.New("db down"), ok: false},
	}
	dto := &dto.IncidentDTO{
		Lat: 1,
		Lon: 1,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(tt.query)).WithArgs(dto.Lat, dto.Lon).WillReturnError(tt.err).WillReturnResult(sqlmock.NewResult(1, 1))
			err := repo.CreateIncidents(t.Context(), dto)
			if tt.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}
func TestUpdateZones(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repo.NewIncidentRepo(db)
	q := `UPDATE zones SET is_dangerous = TRUE WHERE zone_id = (SELECT z.zone_id FROM zones z WHERE ST_Contains(
			z.wkb_geometry,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)
			) LIMIT 1)`
	tests := []struct {
		name  string
		query string
		err   error
		ok    bool
	}{
		{name: "ok", query: q, err: nil, ok: true},
		{name: "error", query: q, err: errors.New("db down"), ok: false},
	}
	dto := &dto.IncidentDTO{
		Lat: 1,
		Lon: 1,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(tt.query)).WithArgs(dto.Lat, dto.Lon).WillReturnError(tt.err).WillReturnResult(sqlmock.NewResult(1, 1))
			err := repo.UpdateZones(t.Context(), dto)
			if tt.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}
func TestUpdateIncidentsByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repo.NewIncidentRepo(db)
	q := `UPDATE incidents SET lat = $1, lon = $2, update_at = NOW() WHERE incident_id = $3`
	tests := []struct {
		name  string
		query string
		err   error
		ok    bool
	}{
		{name: "ok", query: q, err: nil, ok: true},
		{name: "error", query: q, err: errors.New("db down"), ok: false},
	}
	dto := &dto.IncidentDTO{
		Lat: 1,
		Lon: 1,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(tt.query)).WithArgs(dto.Lat, dto.Lon, id).WillReturnError(tt.err).WillReturnResult(sqlmock.NewResult(1, 1))
			err := repo.UpdateIncidentsByID(t.Context(), dto, id)
			if tt.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}
func TestDeleteIncidentsByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repo.NewIncidentRepo(db)
	q := `UPDATE incidents SET status = 'hide', deleted_at = NOW() WHERE incident_id = $1`
	tests := []struct {
		name  string
		query string
		err   error
		ok    bool
	}{
		{name: "ok", query: q, err: nil, ok: true},
		{name: "error", query: q, err: errors.New("db down"), ok: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(tt.query)).WithArgs(id).WillReturnError(tt.err).WillReturnResult(sqlmock.NewResult(1, 1))
			err := repo.DeleteIncidentsByID(t.Context(), id)
			if tt.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}
