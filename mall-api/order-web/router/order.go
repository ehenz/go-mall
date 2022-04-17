package router

import (
	"mall-api/order-web/api/order"
	"mall-api/order-web/middleware"

	"github.com/gin-gonic/gin"
)

func InitOrderRouter(Router *gin.RouterGroup) {
	OrderRouter := Router.Group("order").Use(middleware.JWTAuth())
	{
		OrderRouter.GET("", order.List)       // 订单列表
		OrderRouter.GET("/:id", order.Detail) // 订单详情
		// TODO OrderRouter.DELETE("/:id", order.Delete) // 删除订单
		OrderRouter.POST("", order.New) // 新建订单
		// TODO OrderRouter.PUT("/:id", order.Update)    // 修改订单
	}
}
