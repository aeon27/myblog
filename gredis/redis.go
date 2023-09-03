package gredis

import (
	"encoding/json"
	"time"

	"github.com/aeon27/myblog/pkg/setting"
	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

// 初始化redis相关连接和设置
func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,

		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := conn.Do("AUTH", setting.RedisSetting.Password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, err
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}

	return nil
}

// 设置redis键值对
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// 获取redis键对应的值的字节数组
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return value, err
}

// 判断redis键是否存在
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// 删除redis键
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// 模糊删除redis键
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, k := range keys {
		_, err := conn.Do("DEL", k)
		if err != nil {
			return err
		}
	}

	return nil
}
