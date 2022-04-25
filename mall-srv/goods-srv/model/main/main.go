package main

import (
	"context"
	"log"
	"mall-srv/goods-srv/global"
	"mall-srv/goods-srv/initialize"
	"mall-srv/goods-srv/model"
	"os"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// InitMySQL 生成MySQL表结构
func InitMySQL() {
	dsn := "root:root@tcp(106.13.213.235:3306)/mshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	_ = db.AutoMigrate(
		&model.Goods{},
	)
}

func SyncMysqlToEs() {
	host := "http://106.13.214.17:9200"
	l := log.New(os.Stdout, "es", log.LstdFlags)
	EsClient, err := elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false), elastic.SetTraceLog(l))
	if err != nil {
		panic(err)
	}

	initialize.InitDB()

	var goods []model.Goods
	global.DB.Find(&goods)
	for _, v := range goods {
		esgoods := model.EsGoods{
			ID:          v.ID,
			CategoryID:  v.CategoryId,
			BrandsID:    v.BrandId,
			OnSale:      v.OnSale,
			ShipFree:    v.ShipFree,
			IsNew:       v.IsNew,
			IsHot:       v.IsHot,
			Name:        v.Name,
			ClickNum:    v.ClickNum,
			SoldNum:     v.SoldNum,
			FavNum:      v.FavNum,
			MarketPrice: v.MarketPrice,
			GoodsBrief:  v.GoodsBrief,
			ShopPrice:   v.ShopPrice,
		}

		_, err := EsClient.Index().Index(model.EsGoods{}.GetIndexName()).BodyJson(esgoods).Id(strconv.Itoa(int(esgoods.ID))).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	initialize.InitConfig(true)
	SyncMysqlToEs()
}
