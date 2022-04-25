package global

import (
	"mall-api/userop-web/config"

	pb "mall-api/userop-web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	NacosConfig      *config.NacosConfig = &config.NacosConfig{}
	SrvConfig        *config.SrvConfig   = &config.SrvConfig{}
	Trans            ut.Translator
	GoodsSrvClient   pb.GoodsClient
	AddressSrvClient pb.AddressClient
	MessageSrvClient pb.MessageClient
	UserFavSrvClient pb.UserFavClient
)
