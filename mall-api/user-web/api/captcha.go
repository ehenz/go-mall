package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := cp.Generate()
	if err != nil {
		zap.S().Error("验证码生成错误：", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "验证码生成错误",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId": id,
		"pic":       b64s,
	})
}
