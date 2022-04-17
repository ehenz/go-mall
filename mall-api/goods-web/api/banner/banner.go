package banner

import (
	"context"
	"mall-api/goods-web/api"
	"mall-api/goods-web/forms"
	"mall-api/goods-web/global"
	"mall-api/goods-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
)

func List(c *gin.Context) {
	rsp, err := global.GoodsSrvClient.BannerList(context.Background(), &empty.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	result := make([]interface{}, 0)
	for _, v := range rsp.Data {
		rMap := make(map[string]interface{})
		rMap["id"] = v.Id
		rMap["index"] = v.Index
		rMap["image"] = v.Image
		rMap["url"] = v.Url

		result = append(result, rMap)
	}

	c.JSON(http.StatusOK, result)
}

func New(c *gin.Context) {
	bannerForm := forms.BannerForm{}
	if err := c.ShouldBindJSON(&bannerForm); err != nil {
		api.HandleValidationError(c, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateBanner(context.Background(), &proto.BannerRequest{
		Index: int32(bannerForm.Index),
		Url:   bannerForm.Url,
		Image: bannerForm.Image,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	response := make(map[string]interface{})
	response["id"] = rsp.Id
	response["index"] = rsp.Index
	response["url"] = rsp.Url
	response["image"] = rsp.Image

	c.JSON(http.StatusOK, response)
}

func Update(c *gin.Context) {
	bannerForm := forms.BannerForm{}
	if err := c.ShouldBindJSON(&bannerForm); err != nil {
		api.HandleValidationError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateBanner(context.Background(), &proto.BannerRequest{
		Id:    int32(i),
		Index: int32(bannerForm.Index),
		Url:   bannerForm.Url,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	_, err = global.GoodsSrvClient.DeleteBanner(context.Background(), &proto.BannerRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, "")
}
