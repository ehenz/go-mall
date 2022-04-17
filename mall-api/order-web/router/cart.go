package router

import (
	"mall-api/order-web/api/cart"
	"mall-api/order-web/middleware"

	"github.com/gin-gonic/gin"
)

func InitCartRouter(Router *gin.RouterGroup) {
	// middleware.JWTAuth()
	CartRouter := Router.Group("cart").Use(middleware.JWTAuth())
	{
		CartRouter.GET("", cart.List)          // 购物车详情
		CartRouter.POST("", cart.New)          // 新建
		CartRouter.DELETE("/:id", cart.Delete) // 购物车订单
		CartRouter.PUT("/:id", cart.Update)    // 修改
	}
}
