package conf

import (
	"time"
	"fmt"
	"os"

	"github.com/garyburd/redigo/redis"
)

// redis
func NewRedisPool(pool_cfg *RedisConf) *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:           pool_cfg.MaxIdle,
		MaxActive:         pool_cfg.MaxActive,
		IdleTimeout:       time.Duration(pool_cfg.IdleTimeout) * time.Second,
		HeartbeatMin:      pool_cfg.HeartBeatMin,
		HeartbeatMax:      pool_cfg.HeartBeatMax,
		HeartbeatInterval: time.Duration(pool_cfg.HeartBeatInterval) * time.Second,
		Wait:              pool_cfg.MaxWait,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", pool_cfg.ServerAddr,
				redis.DialConnectTimeout(time.Duration(pool_cfg.ConnTimeout)*time.Millisecond),
				redis.DialReadTimeout(time.Duration(pool_cfg.ReadTimeout)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(pool_cfg.WriteTimeout)*time.Millisecond))
			if err != nil {
				return nil, err
			}
			if len(pool_cfg.Auth) > 0 {
				if _, err := c.Do("AUTH", pool_cfg.Auth); err != nil {
					fmt.Println("auth err:", err.Error())
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	c := pool.Get()
	defer c.Close()
	if _, err := c.Do("ping"); err != nil {
		fmt.Printf("redis init failed, err=%s\n", err.Error())
		os.Exit(1)
	} else {
		fmt.Printf("redis init success conf_name: %s, conf_addr: %s\n", pool_cfg.ConfName, pool_cfg.ServerAddr)
	}

	return pool
}