package category_brand

import (
	"context"
	"mall-api/goods-web/api"
	"mall-api/goods-web/forms"
	"mall-api/goods-web/global"
	"mall-api/goods-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// List 返回所有分类及品牌列表
func List(c *gin.Context) {
	// 所有的list返回的数据结构
	/*
		{
			"total": 100,
			"data":[{},{}]
		}
	*/
	rsp, err := global.GoodsSrvClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	rMap := map[string]interface{}{
		"total": rsp.Total,
	}

	result := make([]interface{}, 0)
	for _, v := range rsp.Data {
		rMap := make(map[string]interface{})
		rMap["id"] = v.Id
		rMap["category"] = map[string]interface{}{
			"id":   v.Category.Id,
			"name": v.Category.Name,
		}
		rMap["brand"] = map[string]interface{}{
			"id":   v.Brand.Id,
			"name": v.Brand.Name,
			"logo": v.Brand.Logo,
		}

		result = append(result, rMap)
	}

	rMap["data"] = result
	c.JSON(http.StatusOK, rMap)
}

// BrandList 根据分类获取品牌列表
func BrandList(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	rsp, err := global.GoodsSrvClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	result := make([]interface{}, 0)
	for _, v := range rsp.Data {
		rMap := make(map[string]interface{})
		rMap["id"] = v.Id
		rMap["name"] = v.Name
		rMap["logo"] = v.Logo

		result = append(result, rMap)
	}

	c.JSON(http.StatusOK, result)
}

func New(c *gin.Context) {
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := c.ShouldBindJSON(&categoryBrandForm); err != nil {
		api.HandleValidationError(c, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	response := make(map[string]interface{})
	response["id"] = rsp.Id

	c.JSON(http.StatusOK, response)
}

func Update(c *gin.Context) {
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := c.ShouldBindJSON(&categoryBrandForm); err != nil {
		api.HandleValidationError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.UpdateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id:         int32(i),
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
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
	_, err = global.GoodsSrvClient.DeleteCategoryBrand(context.Background(), &proto.CategoryBrandRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, "")
}
