package initialize

import (
	"context"
	"fmt"
	"log"
	"mall-srv/goods-srv/global"
	"mall-srv/goods-srv/model"
	"os"

	"go.uber.org/zap"

	"github.com/olivere/elastic/v7"
)

func InitEs() {
	c := global.SrvConfig.EsConfig
	host := fmt.Sprintf("http://%s:%d", c.Host, c.Port)
	logger := log.New(os.Stdout, "es", log.LstdFlags)
	var err error
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}

	// 初始化es的mapping结构
	exist, err := global.EsClient.IndexExists(model.EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exist {
		rsp, err := global.EsClient.CreateIndex(model.EsGoods{}.GetIndexName()).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err != nil {
			panic(err)
		}
		if rsp.Acknowledged {
			zap.S().Info("Es初始化成功")
		} else {
			zap.S().Info("Es初始化失败")
		}
	} else {
		zap.S().Info("Es结构已初始化")
	}
}
