package category

import (
	"context"
	"encoding/json"
	"mall-api/goods-web/api"
	"mall-api/goods-web/forms"
	"mall-api/goods-web/global"
	"mall-api/goods-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
)

func List(c *gin.Context) {
	r, err := global.GoodsSrvClient.GetAllCategoryList(context.Background(), &empty.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
	}

	data := make([]interface{}, 0)
	err = json.Unmarshal([]byte(r.JsonData), &data)
	if err != nil {
		zap.S().Errorw("[List] 查询分类列表失败： ", err.Error())
	}

	c.JSON(http.StatusOK, data)
}

func Detail(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	rMap := make(map[string]interface{})
	subCategory := make([]interface{}, 0)
	if r, err := global.GoodsSrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: int32(i),
	}); err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	} else {
		//写文档 特别是数据多的时候很慢， 先开发后写文档
		for _, value := range r.SubCategorys {
			subCategory = append(subCategory, map[string]interface{}{
				"id":              value.Id,
				"name":            value.Name,
				"level":           value.Level,
				"parent_category": value.ParentCategory,
				"is_tab":          value.IsTab,
			})
		}
		rMap["id"] = r.Info.Id
		rMap["name"] = r.Info.Name
		rMap["level"] = r.Info.Level
		rMap["parent_category"] = r.Info.ParentCategory
		rMap["is_tab"] = r.Info.IsTab
		rMap["sub_categorys"] = subCategory

		c.JSON(http.StatusOK, rMap)
	}
	return
}

func New(c *gin.Context) {
	categoryForm := forms.CategoryForm{}
	if err := c.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidationError(c, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           categoryForm.Name,
		IsTab:          *categoryForm.IsTab,
		Level:          categoryForm.Level,
		ParentCategory: categoryForm.ParentCategory,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id
	request["name"] = rsp.Name
	request["parent"] = rsp.ParentCategory
	request["level"] = rsp.Level
	request["is_tab"] = rsp.IsTab

	c.JSON(http.StatusOK, request)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// 先查询出该分类写的所有子分类
	// 将所有的分类全部逻辑删除
	// 将该分类下的所有的商品逻辑删除（可选）
	_, err = global.GoodsSrvClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func Update(c *gin.Context) {
	categoryForm := forms.UpdateCategoryForm{}
	if err := c.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidationError(c, err)
		return
	}

	id := c.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	request := &proto.CategoryInfoRequest{
		Id:   int32(i),
		Name: categoryForm.Name,
	}
	if categoryForm.IsTab != nil {
		request.IsTab = *categoryForm.IsTab
	}
	_, err = global.GoodsSrvClient.UpdateCategory(context.Background(), request)
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	c.Status(http.StatusOK)
}
