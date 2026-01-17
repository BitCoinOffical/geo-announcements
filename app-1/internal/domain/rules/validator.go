package rules

import (
	"github.com/go-playground/validator/v10"
)

func ValidateLat(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val >= -90 && val <= 90
}

func ValidateLon(fl validator.FieldLevel) bool {
	val := fl.Field().Float()
	return val >= -180 && val <= 180
}
