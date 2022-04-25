package main

import (
	"flag"
	"fmt"
	"mall-srv/stock-srv/global"
	"mall-srv/stock-srv/handler"
	"mall-srv/stock-srv/initialize"
	"mall-srv/stock-srv/proto"
	"mall-srv/stock-srv/utils"
	"mall-srv/stock-srv/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nacos-group/nacos-sdk-go/inner/uuid"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	c := global.SrvConfig

	debug := flag.Bool("debug", true, "是否启用debug模式 true / false")
	flag.Parse()

	initialize.InitLogger()       // 初始化日志
	initialize.InitConfig(*debug) // 初始化配置
	initialize.InitDB()           // 初始化数据库
	initialize.InitRedisSync()    // 初始化redis分布式锁

	// 动态分配一个可用端口
	var err error
	global.SrvConfig.Port, err = utils.GetFreePort()
	if err != nil {
		zap.S().Error("[GetFreePort] 获取可用端口失败")
	}

	// global.SrvConfig.Port = 9999 // 测试时使用 - 固定端口号

	zap.S().Infof(" Listening and serving HTTP on %s:%d\n", global.SrvConfig.Host, global.SrvConfig.Port)

	server := grpc.NewServer()
	// 注册用户服务
	proto.RegisterStockServer(server, &handler.StockServer{})
	// 注册 consul 官方的rpc健康检查服务
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 设置服务监听地址
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.SrvConfig.Host, global.SrvConfig.Port))
	if err != nil {
		panic("fail to listen:" + err.Error())
	}

	// 启动服务
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("fail to serve:" + err.Error())
		}
	}()

	// 监听rocketmq的库存归还请求
	initialize.StockRollbackConsumer()

	// 服务注册到 consul
	serviceUuid, _ := uuid.NewV4()
	serviceId := fmt.Sprintf("%s", serviceUuid)
	registerClient := consul.NewRegisterClient(c.ConsulConfig.Host, c.ConsulConfig.Port)
	registerClient.Regis(c.Name, c.Host, c.Port, c.Tags, serviceId)

	// 优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	registerClient.DeRegis(serviceId)
}
