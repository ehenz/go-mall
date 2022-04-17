package initialize

import (
	"mall-api/order-web/middleware"
	"mall-api/order-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	// 配置跨域
	Router.Use(middleware.Cors())

	ApiGroup := Router.Group("/o/v1")
	router.InitOrderRouter(ApiGroup)
	router.InitCartRouter(ApiGroup)
	return Router
}
