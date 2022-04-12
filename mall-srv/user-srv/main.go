package main

import (
	"flag"
	"fmt"
	"mall-srv/user-srv/global"
	"mall-srv/user-srv/handler"
	"mall-srv/user-srv/initialize"
	"mall-srv/user-srv/proto"
	"mall-srv/user-srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	debug := flag.Bool("debug", true, "是否启用debug模式 true / false")
	flag.Parse()

	initialize.InitLogger()
	initialize.InitConfig(*debug)
	initialize.InitDB()

	// 动态分配一个可用端口
	var err error
	global.SrvConfig.Port, err = utils.GetFreePort()
	if err != nil {
		zap.S().Error("[GetFreePort] 获取可用端口失败")
	}

	zap.S().Infof(" Listening and serving HTTP on %s:%d\n", global.SrvConfig.Host, global.SrvConfig.Port)

	server := grpc.NewServer()
	// 注册用户服务
	proto.RegisterUserServer(server, &handler.UserServer{})
	// 注册 consul 官方的rpc健康检查服务
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.SrvConfig.ConsulConfig.Host, global.SrvConfig.ConsulConfig.Port)
	client, _ := api.NewClient(cfg)
	// 本地测试需要暂时关闭consul的健康检查
	//check := &api.AgentServiceCheck{
	//	GRPC:                           fmt.Sprintf("%s:%d", global.SrvConfig.Host, global.SrvConfig.Port),
	//	Timeout:                        "5s",
	//	Interval:                       "5s",
	//	DeregisterCriticalServiceAfter: "10s",
	//}
	srvId := fmt.Sprintf("%s", uuid.NewV4())
	regis := &api.AgentServiceRegistration{
		Name:    global.SrvConfig.Name,
		ID:      srvId,
		Address: global.SrvConfig.Host,
		Port:    global.SrvConfig.Port,
		Tags:    []string{"user", "srv"},
		// Check:   check,
	}
	_ = client.Agent().ServiceRegister(regis)

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
	// 优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = client.Agent().ServiceDeregister(srvId); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
