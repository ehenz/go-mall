package global

import (
	"mall-api/user-web/config"

	"mall-api/user-web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	NacosConfig   *config.NacosConfig = &config.NacosConfig{}
	SrvConfig     *config.SrvConfig   = &config.SrvConfig{}
	Trans         ut.Translator
	UserSrvClient proto.UserClient
)
