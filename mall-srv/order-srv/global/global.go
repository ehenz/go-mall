package global

import (
	"mall-srv/order-srv/config"
	"mall-srv/order-srv/proto"

	"gorm.io/gorm"
)

var (
	DB             *gorm.DB
	SrvConfig      *config.SrvConfig   = &config.SrvConfig{}
	NacosConfig    *config.NacosConfig = &config.NacosConfig{}
	GoodsSrvClient proto.GoodsClient
	StockSrvClient proto.StockClient
)
