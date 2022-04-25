package main

import (
	"flag"
	"fmt"
	"mall-api/goods-web/global"
	"mall-api/goods-web/initialize"
	"mall-api/goods-web/utils"
	"mall-api/goods-web/utils/register/consul"
	"os"
	"os/signal"
	"syscall"

	"github.com/nacos-group/nacos-sdk-go/inner/uuid"

	"go.uber.org/zap"
)

func main() {
	c := global.SrvConfig

	debug := flag.Bool("debug", true, "是否以debug模式启动")
	if *debug == true {
		c.Port = 8081
	} else {
		c.Port, _ = utils.GetFreePort()
	}

	// 初始化logger
	initialize.InitLogger()
	// 初始化config
	initialize.InitConfig(*debug)
	// 初始化routers
	Routers := initialize.Routers()
	// 初始化表单验证翻译功能
	_ = initialize.InitTrans("zh")
	// 初始化rpc服务客户端
	initialize.InitSrvClient()
	// 初始化Sentinel
	initialize.InitSentinel()

	// 服务注册到 consul
	serviceUuid, _ := uuid.NewV4()
	serviceId := fmt.Sprintf("%s", serviceUuid)
	registerClient := consul.NewRegisterClient(c.ConsulConfig.Host, c.ConsulConfig.Port)
	registerClient.Regis(c.Name, c.Host, c.Port, c.Tags, serviceId)

	go func() {
		err := Routers.Run(fmt.Sprintf("%s:%d", c.Host, c.Port))
		if err != nil {
			zap.S().Error("启动失败：", err.Error())
		}

	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	registerClient.DeRegis(serviceId)
}
