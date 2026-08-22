package redis

import (
	"Blog/settings"
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

//初始化全局redis客户端
var rdb *redis.Client
func Init() error {
	var err error
	rdb  := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s",settings.Cnf.Redis.Host,settings.Cnf.Redis.Port),
		Password: settings.Cnf.Redis.Passwrod,
		DB: 0,
		PoolSize: 10,
	})

	_ , err = rdb.Ping(context.Background()).Result()
	if err != nil {
		return errors.New("redis初始化化失败")
	}
	return nil 
}

func Close() {
	rdb.Close()
}