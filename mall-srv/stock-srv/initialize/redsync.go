package initialize

import (
	"mall-srv/stock-srv/global"

	"github.com/go-redis/redis"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis"
)

func InitRedisSync() {
	// TODO 配置从 nacos 获取
	client := redis.NewClient(&redis.Options{
		Addr: "106.13.213.235:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	global.Rs = redsync.New(pool)
}
