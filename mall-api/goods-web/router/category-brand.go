package router

import (
	cb "mall-api/goods-web/api/category-brand"
	"mall-api/goods-web/middleware"

	"github.com/gin-gonic/gin"
)

func InitCategoryBrandRouter(Router *gin.RouterGroup) {
	CategoryBrandRouter := Router.Group("category_brand").Use(middleware.JaegerTrace())
	{
		CategoryBrandRouter.GET("", cb.List)          // 返回所有分类及品牌列表
		CategoryBrandRouter.DELETE("/:id", cb.Delete) // 删除id类别的品牌
		CategoryBrandRouter.POST("", cb.New)          // 新建类别品牌
		CategoryBrandRouter.PUT("/:id", cb.Update)    // 修改id类别的品牌
		CategoryBrandRouter.GET("/:id", cb.BrandList) // 根据分类获取品牌列表
	}
}
