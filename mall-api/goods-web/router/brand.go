package router

import (
	"mall-api/goods-web/api/brand"
	"mall-api/goods-web/middleware"

	"github.com/gin-gonic/gin"
)

func InitBrandRouter(Router *gin.RouterGroup) {
	BrandRouter := Router.Group("brand").Use(middleware.JaegerTrace())
	{
		BrandRouter.GET("", brand.List)          // 品牌列表
		BrandRouter.DELETE("/:id", brand.Delete) // 删除品牌
		BrandRouter.POST("", brand.New)          // 新建品牌
		BrandRouter.PUT("/:id", brand.Update)    // 修改品牌
	}
}
