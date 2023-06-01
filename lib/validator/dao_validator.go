package validator

import (
	"gin_websocket/dao"

	pValidator "github.com/go-playground/validator/v10"
)

func roleValidator(fl pValidator.FieldLevel) bool {
	roleId, ok := fl.Field().Interface().(int)
	if ok && dao.ExistsRole(roleId) {
		return true
	}
	return false
}
