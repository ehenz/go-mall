package consul

import (
	"fmt"
	"mall-srv/order-srv/global"

	"go.uber.org/zap"

	"github.com/hashicorp/consul/api"
)

type RegisterClient interface {
	Regis(name string, host string, port int, tags []string, id string)
	DeRegis(serviceId string)
}

type Register struct {
	Host string
	Port int
}

func NewRegisterClient(host string, port int) RegisterClient {
	return &Register{
		Host: host,
		Port: port,
	}
}

func (r *Register) Regis(name string, host string, port int, tags []string, id string) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.SrvConfig.ConsulConfig.Host, global.SrvConfig.ConsulConfig.Port)
	client, _ := api.NewClient(cfg)
	// TODO 本地测试需要暂时关闭consul的健康检查
	//check := &api.AgentServiceCheck{
	//	GRPC:                           fmt.Sprintf("%s:%d", global.SrvConfig.Host, global.SrvConfig.Port),
	//	Timeout:                        "5s",
	//	Interval:                       "5s",
	//	DeregisterCriticalServiceAfter: "10s",
	//}
	regis := &api.AgentServiceRegistration{
		Name:    name,
		ID:      id,
		Address: host,
		Port:    port,
		Tags:    tags,
		// Check:   check,
	}
	err := client.Agent().ServiceRegister(regis)
	if err != nil {
		zap.S().Errorw("服务注册失败：", err.Error())
	}
	zap.S().Info("服务注册成功")
}

func (r *Register) DeRegis(serviceId string) {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.SrvConfig.ConsulConfig.Host, global.SrvConfig.ConsulConfig.Port)
	client, _ := api.NewClient(cfg)

	err := client.Agent().ServiceDeregister(serviceId)
	if err != nil {
		zap.S().Errorw("服务注销失败")
	}
	zap.S().Info("服务注销成功")
}
