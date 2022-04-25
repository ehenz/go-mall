package notify

import (
	"context"
	"fmt"
	"mall-api/order-web/api"
	"mall-api/order-web/global"
	pb "mall-api/order-web/proto"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/smartwalle/alipay/v3"
)

func Notify(ctx *gin.Context) {
	c := global.SrvConfig.Alipay
	client, err := alipay.New(c.Appid, c.PrivateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(c.AliPublicKey)
	if err != nil {
		panic(err)
	}

	var notify, _ = client.GetTradeNotification(ctx.Request)
	if notify != nil {
		fmt.Println("交易状态为:", notify.TradeStatus)
		return
	}

	_, err = global.OrderSrvClient.UpdateOrder(context.Background(), &pb.OrderStatus{
		OrderSn: notify.OutTradeNo,
		Status:  string(notify.TradeStatus),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(ctx, err)
		return
	}

	ctx.String(http.StatusOK, "success")
}
