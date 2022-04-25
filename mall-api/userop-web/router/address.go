package router

import (
	"mall-api/userop-web/api/address"
	"mall-api/userop-web/middleware"

	"github.com/gin-gonic/gin"
)

func InitAddressRouter(Router *gin.RouterGroup) {
	AddressRouter := Router.Group("address").Use(middleware.JWTAuth())
	{
		AddressRouter.GET("", address.List)
		AddressRouter.DELETE("/:id", address.Delete)
		AddressRouter.POST("", address.New)
		AddressRouter.PUT("/:id", address.Update)
	}
}
