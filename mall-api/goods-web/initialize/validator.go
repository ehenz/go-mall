package initialize

import (
	"fmt"
	"mall-api/goods-web/global"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	zh_translation "github.com/go-playground/validator/v10/translations/zh"
)

// InitTrans 表单验证信息转中文（官方支持）
func InitTrans(locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New()
		enT := en.New()
		// 第一个参数是备用的语言环境
		uni := ut.New(enT, zhT, enT)
		global.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s)", locale)
		}
		switch locale {
		case "en":
			_ = en_translation.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			_ = zh_translation.RegisterDefaultTranslations(v, global.Trans)
		default:
			_ = en_translation.RegisterDefaultTranslations(v, global.Trans)
		}
		return
	}
	return nil
}
