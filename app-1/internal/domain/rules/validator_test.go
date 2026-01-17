package rules_test

import (
	"testing"

	"github.com/BitCoinOffical/geo-announcements/app-1/internal/domain/rules"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestDTO struct {
	Lat float64 `validate:"lat"`
	Lon float64 `validate:"lon"`
}

func TestValidate_ok(t *testing.T) {
	v := validator.New()

	err := v.RegisterValidation("lat", rules.ValidateLat)
	require.NoError(t, err)

	err = v.RegisterValidation("lon", rules.ValidateLon)
	require.NoError(t, err)

	tests := []struct {
		name string
		dto  TestDTO
		ok   bool
	}{
		{"min ok", TestDTO{Lat: -90, Lon: -180}, true},
		{"max ok", TestDTO{Lat: 90, Lon: 180}, true},
		{"zero ok", TestDTO{Lat: 0, Lon: 0}, true},
		{"too low", TestDTO{Lat: -91, Lon: -181}, false},
		{"too high", TestDTO{Lat: 91, Lon: 181}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.dto)
			if tt.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
