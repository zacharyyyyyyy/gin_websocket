package validator

import (
	"errors"
	"fmt"
	"gin_websocket/lib/logger"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New()
		uni := ut.New(zhT, zhT)
		trans, _ = uni.GetTranslator("zh")
		// 注册翻译器
		zhTranslations.RegisterDefaultTranslations(v, trans)
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
			if name == "-" {
				name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
				if name == "-" {
					return ""
				}
			}
			return name
		})
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("existsAdminRole", roleValidator)
		v.RegisterTranslation("existsAdminRole",
			trans,
			registerTranslator("existsAdminRole", "必须为存在角色"),
			translate,
		)
	}
}

func GetValidMsg(err error) string {
	if errors.Is(err, strconv.ErrSyntax) {
		//暂时未处理参数类型不一致的情况
		return fmt.Sprintf("存在参数类型错误")
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	} else {
		var errString string
		for key, val := range errs.Translate(trans) {
			fmt.Println(key, val)
			errString = val
			break
		}
		return errString
	}
}

func registerTranslator(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

// translate 自定义字段的翻译方法
func translate(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		logger.Runtime.Error("翻译错误:" + err.Error())
		return "存在未定义错误"
	}
	return msg
}
