package global

import (
	"io"
	"mall-srv/order-srv/config"
	"mall-srv/order-srv/proto"

	"github.com/opentracing/opentracing-go"

	"gorm.io/gorm"
)

var (
	DB             *gorm.DB
	SrvConfig      *config.SrvConfig   = &config.SrvConfig{}
	NacosConfig    *config.NacosConfig = &config.NacosConfig{}
	GoodsSrvClient proto.GoodsClient
	StockSrvClient proto.StockClient
	Tracer         opentracing.Tracer
	Closer         io.Closer
)
