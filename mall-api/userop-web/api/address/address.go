package address

import (
	"context"
	"mall-api/userop-web/api"
	"mall-api/userop-web/forms"
	"mall-api/userop-web/global"
	"mall-api/userop-web/model"
	"mall-api/userop-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(c *gin.Context) {
	request := &proto.AddressRequest{}

	claims, _ := c.Get("claims")
	currentUser := claims.(*model.CustomClaims)

	if currentUser.AuthorityID != 1 {
		userId, _ := c.Get("userId")
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.AddressSrvClient.GetAddressList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取地址列表失败")
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["user_id"] = value.UserId
		reMap["province"] = value.Province
		reMap["city"] = value.City
		reMap["district"] = value.District
		reMap["address"] = value.Address
		reMap["signer_name"] = value.SignerName
		reMap["signer_mobile"] = value.SignerMobile

		result = append(result, reMap)
	}

	reMap["data"] = result

	c.JSON(http.StatusOK, reMap)
}

func New(c *gin.Context) {
	addressForm := forms.AddressForm{}
	if err := c.ShouldBindJSON(&addressForm); err != nil {
		api.HandleValidationError(c, err)
		return
	}

	userId, _ := c.Get("userId")
	rsp, err := global.AddressSrvClient.CreateAddress(context.Background(), &proto.AddressRequest{
		UserId:       int32(userId.(uint)),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})

	if err != nil {
		zap.S().Errorw("新建地址失败")
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	_, err = global.AddressSrvClient.DeleteAddress(context.Background(), &proto.AddressRequest{Id: int32(i)})
	if err != nil {
		zap.S().Errorw("删除地址失败")
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func Update(c *gin.Context) {
	addressForm := forms.AddressForm{}
	if err := c.ShouldBindJSON(&addressForm); err != nil {
		api.HandleValidationError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	_, err = global.AddressSrvClient.UpdateAddress(context.Background(), &proto.AddressRequest{
		Id:           int32(i),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})
	if err != nil {
		zap.S().Errorw("更新地址失败")
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
