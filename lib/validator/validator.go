package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func GetValidMsg(err error, obj interface{}) string {
	fmt.Println(&obj)
	getObj := reflect.TypeOf(obj)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			if f, exist := getObj.Elem().FieldByName(e.Field()); exist {
				return f.Tag.Get("msg")
			}
		}
	}
	if errors.Is(err, strconv.ErrSyntax) {
		//暂时未处理参数类型不一致的情况
		return fmt.Sprintf("存在参数类型错误")
	}
	return err.Error()
}