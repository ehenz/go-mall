package initialize

import (
	"encoding/json"
	"mall-api/user-web/global"

	"go.uber.org/zap"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/spf13/viper"
)

// InitConfig 使用 nacos 作为配置中心进行初始化
func InitConfig(debug bool) {
	// 从本地文件读取 nacos 配置信息
	var configFileName string
	configFileName = "user-web/nacos-dev.yaml"
	if debug == false {
		configFileName = "user-web/nacos-pro.yaml"
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	_ = v.ReadInConfig()
	err := v.Unmarshal(global.NacosConfig)
	if err != nil {
		panic(err)
	}

	// 拉取 nacos 配置信息
	c := global.NacosConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         c.Namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: c.Host,
			Port:   c.Port,
		},
	}

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		zap.S().Fatal("[NewConfigClient], 初始化nacos客户端失败")
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: c.DataId,
		Group:  c.Group,
	})
	if err != nil {
		zap.S().Fatal("[GetConfig], 获取nacos配置信息失败")
	}

	//待完成：监听变化并重启服务
	//err = configClient.ListenConfig(vo.ConfigParam{
	//	DataId: c.DataId,
	//	Group:  c.Group,
	//	OnChange: func(namespace, group, dataId, data string) {
	//		zap.S().Info("nacos配置发生变化")
	//		err = json.Unmarshal([]byte(data), global.SrvConfig)
	//		if err != nil {
	//			zap.S().Fatal("[json.Unmarshal]，使用nacos初始化SrvConfig失败", err.Error())
	//		}
	//		InitSrvClient()
	//		// fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
	//	},
	//})
	//if err != nil {
	//	zap.S().Fatal("[ListenConfig], 监听nacos配置信息失败")
	//}

	// 使用拉取的信息初始化 SrvConfig
	err = json.Unmarshal([]byte(content), global.SrvConfig)
	if err != nil {
		zap.S().Fatal("[json.Unmarshal]，使用nacos初始化SrvConfig失败", err.Error())
	}
}

// InitConfig 原始本地Config初始化
//func InitConfig(debug bool) {
//	var configFileName string
//	configFileName = "user-web/config-debug.yaml"
//	if debug == false {
//		configFileName = "user-web/config-pro.yaml"
//	}
//	v := viper.New()
//	v.SetConfigFile(configFileName)
//	_ = v.ReadInConfig()
//	err := v.Unmarshal(global.SrvConfig)
//	if err != nil {
//		panic(err)
//	}
//	zap.S().Infof("配置信息：%v", global.SrvConfig)
//
//	// 监控变化
//	v.WatchConfig()
//	v.OnConfigChange(func(e fsnotify.Event) {
//		zap.S().Infof("配置文件 %v 产生变化", e.Name)
//		_ = v.ReadInConfig()
//		_ = v.Unmarshal(global.SrvConfig)
//		zap.S().Infof("配置信息：%v", global.SrvConfig)
//	})
//}
