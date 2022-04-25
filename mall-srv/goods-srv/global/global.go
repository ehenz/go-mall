package global

import (
	"mall-srv/goods-srv/config"

	"github.com/olivere/elastic/v7"

	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	SrvConfig   *config.SrvConfig   = &config.SrvConfig{}
	NacosConfig *config.NacosConfig = &config.NacosConfig{}
	EsClient    *elastic.Client
)
