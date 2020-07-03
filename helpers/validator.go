package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/lexkong/log"
	"gopkg.in/go-playground/validator.v9"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
	"time"
)

var (
	uni   *ut.UniversalTranslator
	Trans ut.Translator
)


var formatData validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(string)
	if ok {
		format := fl.Param()
		if format == "" {
			format = "2006-01-02 15:04:05"
		}
		_, err := time.Parse(format, date)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	return true
}

func InitValidator() {
	zhCn := zh.New()
	uni = ut.New(zhCn, zhCn)

	Trans, _ = uni.GetTranslator("zh")
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := zh_translations.RegisterDefaultTranslations(validate, Trans)
		if err != nil {
			log.Fatalf(err, "init validator error. errMsg: %s")
		}
		// 注册自定义验证器
		if err := validate.RegisterValidation("formatData", formatData); err != nil {
			log.Fatalf(err, "register custom validation error. errMsg: %s")
		}
	}

}


