package utils

import (
	"github.com/go-playground/validator/v10"
)
func ValidateStruct(s interface{}) error {
	val := validator.New()

	val.Struct(s)

	return nil
}