package core

import (
	"fmt"
	"sw/global"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() {
	fmt.Println(fmt.Sprintf("%s:%s", global.Config.Redis.Host, global.Config.Redis.Port))
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", global.Config.Redis.Host, global.Config.Redis.Port),
		Password:     global.Config.Redis.Password, // no password set
		DB:           global.Config.Redis.Db,       // use default DB
		PoolSize:     32,
		MinIdleConns: 10,
		WriteTimeout: 10 * time.Second,
	})
	global.Redis = rdb
	_, err := global.Redis.Ping(global.Ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("===========redis连接成功=============")
}
