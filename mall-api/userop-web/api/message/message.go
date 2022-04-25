package message

import (
	"context"
	"mall-api/userop-web/api"
	"mall-api/userop-web/forms"
	"mall-api/userop-web/global"
	"mall-api/userop-web/model"
	"mall-api/userop-web/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(c *gin.Context) {
	request := &proto.MessageRequest{}

	userId, _ := c.Get("userId")
	claims, _ := c.Get("claims")
	customClaims := claims.(*model.CustomClaims)
	if customClaims.AuthorityID == 1 {
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.MessageSrvClient.MessageList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取留言失败")
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	rMap := map[string]interface{}{
		"total": rsp.Total,
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		rMap := make(map[string]interface{})
		rMap["id"] = value.Id
		rMap["user_id"] = value.UserId
		rMap["type"] = value.MessageType
		rMap["subject"] = value.Subject
		rMap["message"] = value.Message
		rMap["file"] = value.File

		result = append(result, rMap)
	}
	rMap["data"] = result

	c.JSON(http.StatusOK, rMap)
}

func New(c *gin.Context) {
	userId, _ := c.Get("userId")

	messageForm := forms.MessageForm{}
	if err := c.ShouldBindJSON(&messageForm); err != nil {
		api.HandleValidationError(c, err)
		return
	}

	rsp, err := global.MessageSrvClient.CreateMessage(context.Background(), &proto.MessageRequest{
		UserId:      int32(userId.(uint)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	})

	if err != nil {
		zap.S().Errorw("添加留言失败")
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}
