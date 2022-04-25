package global

import (
	"mall-srv/userop-srv/config"

	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	SrvConfig   *config.SrvConfig   = &config.SrvConfig{}
	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)
