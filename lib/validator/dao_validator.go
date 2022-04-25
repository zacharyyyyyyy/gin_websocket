package validator

import (
	"strconv"

	"gin_websocket/dao"

	pValidator "github.com/go-playground/validator/v10"
)

func roleValidator(fl pValidator.FieldLevel) bool {
	roleId, ok := fl.Field().Interface().(string)
	roleIdInt, err := strconv.Atoi(roleId)
	if ok && err == nil && dao.ExistsRole(roleIdInt) {
		return true
	}
	return false
}
