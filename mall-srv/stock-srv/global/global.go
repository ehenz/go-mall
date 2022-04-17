package global

import (
	"mall-srv/stock-srv/config"

	"github.com/go-redsync/redsync/v4"

	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	SrvConfig   *config.SrvConfig   = &config.SrvConfig{}
	NacosConfig *config.NacosConfig = &config.NacosConfig{}
	Rs          *redsync.Redsync
)
