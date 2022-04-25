package router

import (
	"mall-api/userop-web/api/userfav"
	"mall-api/userop-web/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserFavRouter(Router *gin.RouterGroup) {
	UserFavRouter := Router.Group("userfavs").Use(middleware.JWTAuth())
	{
		UserFavRouter.DELETE("/:id", userfav.Delete) // 删除收藏记录
		UserFavRouter.GET("/:id", userfav.Detail)    // 获取收藏记录
		UserFavRouter.POST("", userfav.New)          // 新建收藏记录
		UserFavRouter.GET("", userfav.List)          // 获取当前用户的收藏
	}
}
