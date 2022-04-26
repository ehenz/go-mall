package main

import (
	"log"
	"mall-srv/stock-srv/model"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 生成表结构
func main() {
	dsn := "username:password@tcp(ip:port)/mshop_stock_srv?charset=utf8mb4&parseTime=True&loc=Local"

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

	//_ = db.AutoMigrate(&model.OrderStatus{})
	db.Create(model.OrderStatus{
		OrderSn: "test_order_sn",
		Status:  1,
		Detail: model.OrderDetailList{struct {
			GoodsId  int32
			GoodsNum int32
		}{843, 1}, {844, 1}},
	})
}
