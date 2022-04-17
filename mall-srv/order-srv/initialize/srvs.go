package initialize

import (
	"fmt"
	"mall-srv/order-srv/global"
	pb "mall-srv/order-srv/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// InitServers 初始化第三方微服务连接的客户端
func InitServers() {
	c := global.SrvConfig
	// 商品服务
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", c.ConsulConfig.Host, c.ConsulConfig.Port, c.GoodsSrvConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 轮询
	)
	if err != nil {
		zap.S().Errorw("[goods-srv] 服务连接失败", "msg", err.Error())
	}

	goodsSrvClient := pb.NewGoodsClient(goodsConn)
	global.GoodsSrvClient = goodsSrvClient

	// 库存服务
	stockConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", c.ConsulConfig.Host, c.ConsulConfig.Port, c.StockSrvConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 轮询
	)
	if err != nil {
		zap.S().Errorw("[stock-srv] 服务连接失败", "msg", err.Error())
	}
	stockSrvClient := pb.NewStockClient(stockConn)
	global.StockSrvClient = stockSrvClient

}
