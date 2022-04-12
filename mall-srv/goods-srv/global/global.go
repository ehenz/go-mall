package global

import (
	"mall-srv/goods-srv/config"

	"gorm.io/gorm"
)

var (
	DB          *gorm.DB
	SrvConfig   *config.SrvConfig   = &config.SrvConfig{}
	NacosConfig *config.NacosConfig = &config.NacosConfig{}
)

//func init() {
//	dsn := "root:root@tcp(106.13.213.235:3306)/mshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
//
//	newLogger := logger.New(
//		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
//		logger.Config{
//			SlowThreshold:             time.Second, // 慢 SQL 阈值
//			LogLevel:                  logger.Info, // 日志级别
//			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
//			Colorful:                  true,        // 彩色打印
//		},
//	)
//	var err error
//	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		Logger: newLogger,
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//	})
//	if err != nil {
//		panic(err)
//	}
//}
