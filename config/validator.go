package config

import "github.com/go-playground/validator/v10"

// CustomValidator implements echo.Validator
type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
