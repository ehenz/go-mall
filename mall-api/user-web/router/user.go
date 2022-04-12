package router

import (
	"mall-api/user-web/api"
	"mall-api/user-web/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("/user")
	{
		UserRouter.GET("/list", middleware.JWTAuth(), middleware.IsAdminAuth(), api.GetUserList)
		UserRouter.POST("/pwd_login", api.LoginByPassword)
		UserRouter.POST("register", api.Register)
	}
}
