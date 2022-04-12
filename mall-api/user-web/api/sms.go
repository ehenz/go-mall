package api

import (
	"context"
	"fmt"
	"mall-api/user-web/forms"
	"mall-api/user-web/global"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Sms(c *gin.Context) {
	// 表单验证
	SmsForm := forms.SmsForm{}
	if err := c.ShouldBind(&SmsForm); err != nil {
		HandleValidationError(c, err)
	}

	// 发送验证码逻辑
	smsCode := "123456"

	// 保存至 redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.SrvConfig.RedisConfig.Host, global.SrvConfig.RedisConfig.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rdb.Set(context.Background(), SmsForm.Mobile, smsCode, 300*time.Second)

	c.JSON(http.StatusOK, gin.H{
		"msg": "验证码发送成功",
	})
}
