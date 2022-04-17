package initialize

import (
	"fmt"
	"mall-api/goods-web/global"
	pb "mall-api/goods-web/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

// InitSrvClient 加入负载均衡 grpc-consul-resolver
func InitSrvClient() {
	c := global.SrvConfig
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", c.ConsulConfig.Host, c.ConsulConfig.Port, c.GoodsSrvConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 轮询
	)
	if err != nil {
		zap.S().Errorw("[GetGoodsSrvList] 服务连接失败", "msg", err.Error())
	}

	userSrvClient := pb.NewGoodsClient(conn)
	global.GoodsSrvClient = userSrvClient
}

// InitSrvClient 初始版本
//func InitSrvClient() {
//	// 从注册中心获取用户服务信息 userSrvHost 和 userSrvPort
//	userSrvHost, userSrvPort := "", 0
//	cfg := api.DefaultConfig()
//	cfg.Address = fmt.Sprintf("%s:%d", global.SrvConfig.ConsulConfig.Host, global.SrvConfig.ConsulConfig.Port)
//	client, err := api.NewClient(cfg)
//	if err != nil {
//		panic(err)
//	}
//	userSrv, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.SrvConfig.UserSrvConfig.Name))
//	if err != nil {
//		panic(err)
//	}
//	for _, v := range userSrv {
//		userSrvHost = v.Address
//		userSrvPort = v.Port
//		break
//	}
//
//	if userSrvHost == "" {
//		zap.S().Error("[InitSrvClient] 错误，从注册中心获取用户服务信息失败！")
//	}
//
//	// 拨号连接用户grpc服务
//	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
//	if err != nil {
//		zap.S().Errorw("[GetUserList] 服务连接失败", "msg", err.Error())
//	}
//
//	// 初始化rpc服务客户端
//	// 待解决：服务ip或port改动后，需要重新初始化rpc客户端
//	// 待解决：grpc 连接池 或者 负载均衡
//	global.UserSrvClient = pb.NewUserClient(conn)
//}
