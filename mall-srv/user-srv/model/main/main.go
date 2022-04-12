package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mall-srv/user-srv/model"
	"os"
	"time"

	"github.com/anaskhan96/go-password-encoder"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func genMd5(code string) string {
	hash := md5.New()
	_, _ = io.WriteString(hash, code)
	return hex.EncodeToString(hash.Sum(nil))
}

func main() {
	dsn := "root:root@tcp(106.13.213.235:3306)/mshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

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

	// _ = db.AutoMigrate(&model.User{})
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode("admin123", options)
	pw := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	for i := 0; i < 10; i++ {
		user := model.User{
			BaseModel: model.BaseModel{},
			Mobile:    fmt.Sprintf("1824490801%d", i),
			Password:  pw,
			NickName:  fmt.Sprintf("Jyen%d", i),
			Birthday:  nil,
			Gender:    "male",
			Role:      0,
		}
		db.Save(&user)
	}
}
