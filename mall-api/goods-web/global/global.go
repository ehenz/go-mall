package global

import (
	"mall-api/goods-web/config"

	pb "mall-api/goods-web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	NacosConfig    *config.NacosConfig = &config.NacosConfig{}
	SrvConfig      *config.SrvConfig   = &config.SrvConfig{}
	Trans          ut.Translator
	GoodsSrvClient pb.GoodsClient
)
