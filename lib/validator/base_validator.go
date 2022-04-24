package validator

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

func intValidator(fl validator.FieldLevel) bool {
	param, ok := fl.Field().Interface().(string)
	_, err := strconv.Atoi(param)
	if ok && err == nil {
		return true
	}
	return false
}
