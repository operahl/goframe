package myredis

import (

	"github.com/garyburd/redigo/redis"
	"time"
)

var ErrNil = redis.ErrNil

var MainRds *redis.Pool

func InitRedis(host string,auth string,db int,maxIdle int,MaxActive int ,idleTimeout int) {
	MainRds = InitRed(host,auth,db ,maxIdle,MaxActive,idleTimeout)
}

// 初始化 redis 连接池
func InitRed(host string,auth string,db int,maxIdle int,MaxActive int ,idleTimeout int) *redis.Pool {
	idle := time.Duration(idleTimeout)
	redisPool := &redis.Pool{
		MaxIdle:     maxIdle,   //空闲连接数
		MaxActive:   MaxActive, //最大连接数
		IdleTimeout: idle * time.Second,
		//Wait:        true,
		Dial: func() (conn redis.Conn, err error) {
			return redis.Dial(
				"tcp",
				host,
				redis.DialPassword(auth),
				redis.DialDatabase(db),
			)
		},
	}
	return redisPool
}
