package goods

import (
	"context"
	"fmt"
	"mall-api/goods-web/api"
	"mall-api/goods-web/forms"
	pb "mall-api/goods-web/proto"
	"strconv"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"

	"go.uber.org/zap"

	"mall-api/goods-web/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	req := &pb.GoodsFilterRequest{}

	minPrice := c.DefaultQuery("pmin", "0")
	minPriceInt, _ := strconv.Atoi(minPrice)
	req.PriceMin = int32(minPriceInt)

	maxPrice := c.DefaultQuery("pmax", "0")
	maxPriceInt, _ := strconv.Atoi(maxPrice)
	req.PriceMax = int32(maxPriceInt)

	isHot := c.DefaultQuery("ih", "0")
	if isHot == "1" {
		req.IsHot = true
	}

	isNew := c.DefaultQuery("in", "0")
	if isNew == "1" {
		req.IsNew = true
	}

	isTab := c.DefaultQuery("it", "0")
	if isTab == "1" {
		req.IsTab = true
	}

	categoryId := c.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	req.TopCategory = int32(categoryIdInt)

	page := c.DefaultQuery("p", "0")
	pageInt, _ := strconv.Atoi(page)
	req.Pages = int32(pageInt)

	pageNum := c.DefaultQuery("pn", "0")
	pageNumInt, _ := strconv.Atoi(pageNum)
	req.Pages = int32(pageNumInt)

	keywords := c.DefaultQuery("q", "")
	req.KeyWords = keywords

	brandId := c.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	req.Brand = int32(brandIdInt)

	// Sentinel限流GoodsList访问
	e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		// Blocked. We could get the block reason from the BlockError.
		c.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "服务器繁忙，请稍后重试",
		})
	}
	// TODO 更改所有的api接口
	// 传递gin的Ctx，里面有从路由传递的tracer和start_span
	r, err := global.GoodsSrvClient.GoodsList(context.WithValue(context.Background(), "ginCtx", c), req)
	if err != nil {
		zap.S().Errorw("查询商品列表失败")
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	e.Exit()

	rMap := map[string]interface{}{
		"total": r.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, value := range r.Data {
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	rMap["data"] = goodsList
	c.JSON(http.StatusOK, rMap)
}

func New(c *gin.Context) {
	goodsForm := forms.GoodsForm{}
	if err := c.ShouldBind(&goodsForm); err != nil {
		api.HandleValidationError(c, err)
	}

	rsp, err := global.GoodsSrvClient.CreateGoods(context.Background(), &pb.CreateGoodsInfo{
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		GoodsDesc:       goodsForm.GoodsDesc,
		ShipFree:        goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
	}
	// TODO 结合商品库存服务
	c.JSON(http.StatusOK, rsp)
}

func Detail(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	r, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &pb.GoodInfoRequest{Id: int32(idInt)})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
	}

	rsp := map[string]interface{}{
		"id":          r.Id,
		"name":        r.Name,
		"goods_brief": r.GoodsBrief,
		"desc":        r.GoodsDesc,
		"ship_free":   r.ShipFree,
		"images":      r.Images,
		"desc_images": r.DescImages,
		"front_image": r.GoodsFrontImage,
		"shop_price":  r.ShopPrice,
		"category": map[string]interface{}{
			"id":   r.Category.Id,
			"name": r.Category.Name,
		},
		"brand": map[string]interface{}{
			"id":   r.Brand.Id,
			"name": r.Brand.Name,
			"logo": r.Brand.Logo,
		},
		"is_hot":  r.IsHot,
		"is_new":  r.IsNew,
		"on_sale": r.OnSale,
	}

	c.JSON(http.StatusOK, rsp)
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &pb.DeleteGoodsInfo{Id: int32(idInt)})
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func Stock(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	// TODO 库存
	fmt.Println(idInt)
}

func Status(c *gin.Context) {
	goodsStatusForm := forms.GoodsStatusForm{}
	err := c.ShouldBind(&goodsStatusForm)
	if err != nil {
		api.HandleValidationError(c, err)
	}

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	_, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &pb.CreateGoodsInfo{
		Id:     int32(idInt),
		IsHot:  goodsStatusForm.IsHot,
		IsNew:  goodsStatusForm.IsNew,
		OnSale: goodsStatusForm.OnSale,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}

func Update(c *gin.Context) {
	goodsForm := forms.GoodsForm{}

	err := c.ShouldBindJSON(&goodsForm)
	if err != nil {
		api.HandleValidationError(c, err)
	}

	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Status(http.StatusNotFound)
	}

	_, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &pb.CreateGoodsInfo{
		Id:              int32(idInt),
		Name:            goodsForm.Name,
		GoodsSn:         goodsForm.GoodsSn,
		Stocks:          goodsForm.Stocks,
		MarketPrice:     goodsForm.MarketPrice,
		ShopPrice:       goodsForm.ShopPrice,
		GoodsBrief:      goodsForm.GoodsBrief,
		GoodsDesc:       goodsForm.GoodsDesc,
		ShipFree:        goodsForm.ShipFree,
		Images:          goodsForm.Images,
		DescImages:      goodsForm.DescImages,
		GoodsFrontImage: goodsForm.FrontImage,
		CategoryId:      goodsForm.CategoryId,
		BrandId:         goodsForm.Brand,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
