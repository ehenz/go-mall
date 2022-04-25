package main

import (
	"flag"
	"fmt"
	"mall-api/goods-web/global"
	"mall-api/goods-web/initialize"
	"mall-api/goods-web/utils"
	"mall-api/goods-web/utils/register/consul"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nacos-group/nacos-sdk-go/inner/uuid"

	"go.uber.org/zap"
)

func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

func main() {
	c := global.SrvConfig

	// 获取本机ip地址 - 服务器用
	ip, err := GetOutBoundIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip)
	c.Host = "106.13.214.17"

	// 获取一个可用端口
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
