package middleware

import (
	"mall-api/order-web/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		auth := claims.(*model.CustomClaims).AuthorityID

		if auth != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "无操作权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
