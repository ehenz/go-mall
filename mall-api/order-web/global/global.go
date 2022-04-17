package global

import (
	"mall-api/order-web/config"

	pb "mall-api/order-web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	NacosConfig    *config.NacosConfig = &config.NacosConfig{}
	SrvConfig      *config.SrvConfig   = &config.SrvConfig{}
	Trans          ut.Translator
	GoodsSrvClient pb.GoodsClient
	OrderSrvClient pb.OrderClient
	StockSrvClient pb.StockClient
)
