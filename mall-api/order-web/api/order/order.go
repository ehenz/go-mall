package order

import (
	"context"
	"mall-api/order-web/api"
	"mall-api/order-web/forms"
	"mall-api/order-web/global"
	model2 "mall-api/order-web/model"
	pb "mall-api/order-web/proto"
	"mall-api/order-web/utils/pay"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	claims, _ := c.Get("claims")
	model := claims.(*model2.CustomClaims)
	orderReq := pb.OrderFilterRequest{}
	// 是否为管理员 - 1 为管理员
	if model.AuthorityID != 1 {
		// 非管理员则加上用户id进行限制
		userId, _ := c.Get("userId")
		orderReq.UserId = int32(userId.(uint))
	}
	p := c.DefaultQuery("p", "1")
	pInt, _ := strconv.ParseInt(p, 10, 32)
	pNum := c.DefaultQuery("pnum", "10")
	pNumInt, _ := strconv.ParseInt(pNum, 10, 32)
	orderReq.Pages = int32(pInt)
	orderReq.PagePerNums = int32(pNumInt)

	rsp, err := global.OrderSrvClient.OrderList(context.WithValue(context.Background(), "ginCtx", c), &orderReq)
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  rsp.Data,
	})
}
func New(c *gin.Context) {
	req := pb.OrderRequest{}

	// userId
	userId, _ := c.Get("userId")
	req.UserId = int32(userId.(uint))
	zap.S().Info(userId)

	// 其他表单信息
	orderForm := forms.NewOrderForm{}
	if err := c.ShouldBindJSON(&orderForm); err != nil {
		api.HandleValidationError(c, err)
	}

	req.Name = orderForm.Name
	req.Mobile = orderForm.Mobile
	req.Address = orderForm.Address
	req.Post = orderForm.Post

	rsp, err := global.OrderSrvClient.CreateOrder(context.Background(), &req)
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}

	alipayURL := pay.AlipayPagePay(rsp.OrderSn, strconv.FormatFloat(float64(rsp.Total), 'f', 2, 64))
	c.JSON(http.StatusOK, gin.H{
		"orderId":    rsp.Id,
		"alipay_url": alipayURL,
	})

}
func Detail(c *gin.Context) {
	orderReq := pb.OrderRequest{}

	// 订单id
	orderId := c.Param("id")
	orderIdInt, _ := strconv.ParseInt(orderId, 10, 32)
	orderReq.Id = int32(orderIdInt)

	claims, _ := c.Get("claims")
	model := claims.(*model2.CustomClaims)
	// 是否为管理员 - 1 为管理员
	if model.AuthorityID != 1 {
		// 非管理员则加上用户id进行限制
		userId, _ := c.Get("userId")
		orderReq.UserId = int32(userId.(uint))
	}

	rsp, err := global.OrderSrvClient.OrderDetail(context.Background(), &orderReq)
	if err != nil {
		api.HandleGrpcErrorToHttp(c, err)
		return
	}
	// TODO 格式处理
	c.JSON(http.StatusOK, gin.H{
		"OrderInfo": rsp.OrderInfo,
		"Goods":     rsp.Goods,
	})
}
