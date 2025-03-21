package config

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisService represents the client of redis.
type RedisService struct {
	Client *redis.Client
}

func SetupRedis(cfg *Config) (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:      cfg.REDISHOST,
		Password:  cfg.REDIS_PASSWORD, // Add password for authentication
		DB:        0,
		TLSConfig: &tls.Config{},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.New("failed to connect to redis: " + err.Error())
	}

	return &RedisService{
		Client: client,
	}, nil
}

// func SetupRedis(cfg *Config) (*RedisService, error) {
// 	client := redis.NewClient(&redis.Options{
// 		Addr: cfg.REDISHOST,
// 		DB:   0,
// 	})

// 	//fmt.Println("redis-port address", client)

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	_, err := client.Ping(ctx).Result()
// 	if err != nil {
// 		return nil, errors.New("failed to connect to redis")
// 	}

// 	return &RedisService{
// 		Client: client,
// 	}, nil
// }

// SetDataInRedis will set  data in redis with a key and expiration time.
func (r *RedisService) SetDataInRedis(key string, value []byte, expTime time.Duration) error {
	err := r.Client.Set(context.Background(), key, value, expTime).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetFromRedis will help to retrieve the data from redis.
func (r *RedisService) GetFromRedis(key string) (string, error) {
	jsonData, err := r.Client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return jsonData, nil
}

// DeleteFromRedis will help to delet the data from redis
func (r *RedisService) DeleteFromRedis(key string) error {
	ctx := context.Background()
	err := r.Client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil

}
