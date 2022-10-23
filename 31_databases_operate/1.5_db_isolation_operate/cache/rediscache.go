package cache

import (
	"encoding/json"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/pragmaticrveivews/golang-mux-api/entity"
)

type redisCache struct {
	host    string
	db      int // 0-15
	expires time.Duration
}

/// 构造函数
func NewRedisCache(host string, db int, exp time.Duration) *redisCache {
	return &redisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return &redis.NewClient(&redis.Opeions{
		Addr:     cache.host,
		Passwrod: "",
		DB:       cache.db,
	})

}
func (cache *redisCache) Set(key string, value *entity.Post) {
	client := cache.getClient()
	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	client.Set(key, json, cache.expires*time.Second)
}

func (cache *redisCache) Get(key string, value *entity.Post) *entity.Post {
	client := cache.getClient()
	val, err := client.Get(key).Result()
	if err != nil {
		// panic(err)
		return nil
	}
	post := entity.Post{}
	err2 := json.Unmarshal([]byte(val), &post)
	if err2 != nil {
		panic(err2)
	}
	return &post
}
