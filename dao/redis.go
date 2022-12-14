package dao

import (
	"github.com/go-redis/redis"
	"sanHeRecruitment/config"
)

var (
	Redis   *redis.Client
	RedisDF *redis.Client //redis db 1
)

func InitRedis() (err error) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Addr,
		Password: config.RedisConfig.Password, // no password set
		//Password:"123456",
		DB: config.RedisConfig.DB, // use default DB
	})
	_, err = Redis.Ping().Result()
	if err != nil {
		return err
	}

	RedisDF = redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Addr,
		Password: config.RedisConfig.Password, // no password set
		//Password:"123456",
		DB: 1, // use default DB
	})
	_, err = RedisDF.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}
