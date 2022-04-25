package router

import (
	"mall-api/goods-web/api/banner"
	"mall-api/goods-web/middleware"

	"github.com/gin-gonic/gin"
)

func InitBannerRouter(Router *gin.RouterGroup) {
	BannerRouter := Router.Group("banner").Use(middleware.JaegerTrace())
	{
		BannerRouter.GET("", banner.List)                                                          // 轮播图列表页
		BannerRouter.DELETE("/:id", middleware.JWTAuth(), middleware.IsAdminAuth(), banner.Delete) // 删除轮播图
		BannerRouter.POST("", middleware.JWTAuth(), middleware.IsAdminAuth(), banner.New)          //新建轮播图
		BannerRouter.PUT("/:id", middleware.JWTAuth(), middleware.IsAdminAuth(), banner.Update)    //修改轮播图信息
	}
}
