package cart

import (
	"context"
	"mall-api/order-web/api"
	"mall-api/order-web/forms"
	"mall-api/order-web/global"
	pb "mall-api/order-web/proto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	userId, _ := c.Get("userId")
	//userIdInt, _ := strconv.ParseInt(userId, 10, 32)
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &pb.UserInfo{Id: int32(userId.(uint))})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
	}

	var goodsIdList []int32
	for _, v := range rsp.Data {
		goodsIdList = append(goodsIdList, v.GoodsId)
	}

	if len(goodsIdList) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}

	// 根据goodsIdList找商品信息
	goodsInfoRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &pb.BatchGoodsIdInfo{Id: goodsIdList})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
	}

	// 拼凑商品信息、checked 状态、数量
	rspDate := make([]interface{}, 0)
	for _, k := range rsp.Data {
		for _, v := range goodsInfoRsp.Data {
			if k.GoodsId == v.Id {
				tMap := make(map[string]interface{})
				tMap["id"] = v.Id
				tMap["name"] = v.Name
				tMap["shop_price"] = v.ShopPrice
				tMap["market_price"] = v.MarketPrice
				tMap["front_image"] = v.GoodsFrontImage
				tMap["checked"] = k.Checked
				tMap["num"] = k.Nums

				rspDate = append(rspDate, tMap)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  rspDate,
	})

}
func New(c *gin.Context) {
	cartItem := forms.CartItem{}
	if err := c.ShouldBindJSON(&cartItem); err != nil {
		api.HandleValidationError(c, err)
	}

	// 查看商品是否存在
	_, err := global.GoodsSrvClient.GetGoodsDetail(context.Background(), &pb.GoodInfoRequest{Id: cartItem.GoodsId})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
	}

	// 对比库存是否充足
	stockRsp, err := global.StockSrvClient.CheckStock(context.Background(), &pb.StockInfo{
		GoodsId: cartItem.GoodsId,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
	}
	if cartItem.Nums > stockRsp.Stock {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "库存不足",
		})
		return
	}

	userId, _ := c.Get("userId")
	r, err := global.OrderSrvClient.CreateCart(context.Background(), &pb.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: cartItem.GoodsId,
		Nums:    cartItem.Nums,
		Checked: true,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
	}
	c.JSON(http.StatusOK, gin.H{
		"id": r.Id,
	})
}
func Delete(c *gin.Context) {
	goodsId := c.Param("id")
	goodsIdInt, _ := strconv.ParseInt(goodsId, 10, 32)
	userId, _ := c.Get("userId")
	_, err := global.OrderSrvClient.DeleteCart(context.Background(), &pb.CartItemRequest{GoodsId: int32(goodsIdInt), UserId: int32(userId.(uint))})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}
func Update(c *gin.Context) {
	cartUpdate := forms.CartUpdate{}
	err := c.ShouldBindJSON(&cartUpdate)
	if err != nil {
		api.HandleValidationError(c, err)
		return
	}
	goodsId := c.Param("id")
	goodsIdInt, _ := strconv.ParseInt(goodsId, 10, 32)
	userId, _ := c.Get("userId")
	_, err = global.OrderSrvClient.UpdateCart(context.Background(), &pb.CartItemRequest{
		UserId:  int32(userId.(uint)),
		GoodsId: int32(goodsIdInt),
		Nums:    cartUpdate.Nums,
		Checked: cartUpdate.Checked,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "更改成功",
	})
}
