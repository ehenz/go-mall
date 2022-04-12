package router

import (
	"mall-api/user-web/api"

	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
		BaseRouter.POST("sms", api.Sms)
	}
}
