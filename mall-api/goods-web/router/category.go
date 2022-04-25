package router

import (
	"mall-api/goods-web/api/category"
	"mall-api/goods-web/middleware"

	"github.com/gin-gonic/gin"
)

func InitCategoryRouter(Router *gin.RouterGroup) {
	CategoryRouter := Router.Group("category").Use(middleware.JaegerTrace())
	{
		CategoryRouter.GET("", category.List)          // 商品类别列表页
		CategoryRouter.DELETE("/:id", category.Delete) // 删除分类
		CategoryRouter.GET("/:id", category.Detail)    // 获取分类详情
		CategoryRouter.POST("", category.New)          // 新建分类
		CategoryRouter.PUT("/:id", category.Update)    // 修改分类信息
	}
}
