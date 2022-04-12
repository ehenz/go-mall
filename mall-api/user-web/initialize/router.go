package initialize

import (
	"mall-api/user-web/middleware"
	"mall-api/user-web/router"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	// 配置跨域
	Router.Use(middleware.Cors())

	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	return Router
}
