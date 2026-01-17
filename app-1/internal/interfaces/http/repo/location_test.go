package repo_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/dto"
	"github.com/BitCoinOffical/geo-announcements/app-1/internal/interfaces/http/repo"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	userid = "id"
)

func TestCreateLocation(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repo.NewLocationRepo(db)
	q := `
	    INSERT INTO user_checks_loс (user_id, lat, lon, zone_id)
		SELECT $1, $2, $3, z.zone_id
		FROM zones z
		WHERE ST_Within(ST_SetSRID(ST_MakePoint($3, $2), 4326), z.wkb_geometry)
		LIMIT 1 RETURNING location_id;`
	tests := []struct {
		name  string
		query string
		err   error
		ok    bool
	}{
		{name: "ok", query: q, err: nil, ok: true},
		{name: "error", query: q, err: errors.New("db down"), ok: false},
	}
	dto := &dto.LocationDTO{
		Lat: 1,
		Lon: 1,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec(regexp.QuoteMeta(tt.query)).WithArgs(userid, dto.Lat, dto.Lon).WillReturnError(tt.err).WillReturnResult(sqlmock.NewResult(1, 1))
			err := repo.CreateLocation(t.Context(), dto, userid)
			if tt.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}
}

func TestGetDangerZones(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repo.NewLocationRepo(db)
	q := `	SELECT z.zone_id, z.lat, z.lon,
		ST_Distance(
			z.wkb_geometry,
			ST_SetSRID(ST_Point($1, $2), 4326)
		) AS distanc FROM zones z
		WHERE z.is_dangerous = TRUE
		ORDER BY z.wkb_geometry <-> ST_SetSRID(ST_Point($1, $2), 4326) LIMIT 5;`
	tests := []struct {
		name  string
		query string
		err   error
		ok    bool
	}{
		{name: "ok", query: q, err: nil, ok: true},
		{name: "error", query: q, err: errors.New("db down"), ok: false},
	}
	dto := &dto.LocationDTO{
		Lat: 1,
		Lon: 1,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectQuery(regexp.QuoteMeta(tt.query)).WithArgs(dto.Lat, dto.Lon).WillReturnError(tt.err).WillReturnRows(sqlmock.NewRows([]string{
				"zone_id",
				"lat",
				"lon",
				"distant"}))
			rows, err := repo.GetDangerZones(t.Context(), dto)
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
