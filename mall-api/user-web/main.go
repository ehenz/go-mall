package main

import (
	"flag"
	"fmt"
	"mall-api/user-web/global"
	"mall-api/user-web/initialize"
	"mall-api/user-web/utils"
	myvalidator "mall-api/user-web/validator"

	ut "github.com/go-playground/universal-translator"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"
)

func main() {
	debug := flag.Bool("debug", true, "是否以debug模式启动")
	if *debug == true {
		global.SrvConfig.UserApiConfig.Port = 8080
	} else {
		global.SrvConfig.UserApiConfig.Port, _ = utils.GetFreePort()
	}

	// 初始化logger
	initialize.InitLogger()
	// 初始化config
	initialize.InitConfig(*debug)
	// 初始化routers
	Routers := initialize.Routers()
	// 初始化表单验证翻译功能
	_ = initialize.InitTrans("zh")
	// 初始化rpc服务客户端
	initialize.InitSrvClient()

	// 注册验证器及其错误翻译
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 不是正确的手机号！", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	err := Routers.Run(fmt.Sprintf("%s:%d", global.SrvConfig.UserApiConfig.Host, global.SrvConfig.UserApiConfig.Port))
	if err != nil {
		zap.S().Error("启动失败：", err.Error())
	}
}
