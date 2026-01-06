package rules

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateLat(fl validator.FieldLevel) bool {
	domain := fl.Field().String()
	regex := `^[-+]?(?:[1-8]?\d(?:\.\d+)?|90(?:\.0+)?)$`
	match, _ := regexp.MatchString(regex, domain)
	return match
}

func ValidateLon(fl validator.FieldLevel) bool {
	domain := fl.Field().String()
	regex := `^[-+]?(?:180(?:\.0+)?|(?:1[0-7]\d|[1-9]?\d)(?:\.\d+)?)$`
	match, _ := regexp.MatchString(regex, domain)
	return match
}
