package initialize

import (
	"mall-api/userop-web/middleware"
	"mall-api/userop-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	// 配置跨域
	Router.Use(middleware.Cors())

	ApiGroup := Router.Group("/o/v1")
	router.InitMessageRouter(ApiGroup)
	router.InitAddressRouter(ApiGroup)
	router.InitUserFavRouter(ApiGroup)
	return Router
}
