package router

import (
	"mall-api/goods-web/api/goods"
	"mall-api/goods-web/middleware"

	"github.com/gin-gonic/gin"
)

func InitGoodsRouter(Router *gin.RouterGroup) {
	GoodsRouter := Router.Group("goods")
	{
		GoodsRouter.GET("", goods.List)                                                          // 商品列表
		GoodsRouter.POST("", middleware.JWTAuth(), middleware.IsAdminAuth(), goods.New)          // 增加
		GoodsRouter.GET("/:id", goods.Detail)                                                    // 详情
		GoodsRouter.GET("/:id/stock", goods.Stock)                                               // 库存
		GoodsRouter.DELETE("/:id", middleware.JWTAuth(), middleware.IsAdminAuth(), goods.Delete) // 删除
		GoodsRouter.PATCH("/:id", middleware.JWTAuth(), middleware.IsAdminAuth(), goods.Status)  // 状态
		GoodsRouter.PUT("/:id", middleware.JWTAuth(), middleware.IsAdminAuth(), goods.Update)    // 更改
	}
}
