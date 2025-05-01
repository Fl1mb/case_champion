package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rdb *redis.Client
}

func NewCache(opt *redis.Options) (*Cache, error) {
	db := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	if _, err := db.Ping(ctx).Result(); err != nil {
		_ = db.Close() // Закрываем соединение при ошибке
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Cache{rdb: db}, nil
}

func (c *Cache) Close() error {
	return c.rdb.Close()
}

func (c *Cache) HealthCheck(ctx context.Context) error {
	if c == nil || c.rdb == nil {
		return fmt.Errorf("Cache is not init")
	}
	_, err := c.rdb.Ping(ctx).Result()
	return err
}

func (c *Cache) SetData(ctx context.Context, jwt, login string, user_id int32) error {
	if jwt == "" {
		return fmt.Errorf("empty jwt")
	}
	hashKey := jwt
	fields := map[string]interface{}{
		"login":   login,
		"user_id": user_id,
	}

	return c.rdb.HSet(ctx, hashKey, fields).Err()
}

func (c *Cache) SetDeadTime(ctx context.Context, hash string, dl time.Duration) error {
	return c.rdb.Expire(ctx, hash, dl).Err()
}

func (c *Cache) Delete(ctx context.Context, hash string) error {
	return c.rdb.Del(ctx, hash).Err()
}

func (c *Cache) Exists(ctx context.Context, hash string) bool {
	_, err := c.rdb.Get(ctx, hash).Result()
	return err != redis.Nil
}

func (c *Cache) GetData(ctx context.Context, hash string) (int, string, error) {
	if !c.Exists(ctx, hash) {
		return 0, "", fmt.Errorf("non-existent key")
	}
	val, err := c.rdb.HGetAll(ctx, hash).Result()
	if err == redis.Nil {
		return 0, "", err
	} else if err != nil {
		return 0, "", fmt.Errorf("Ошибка чтения хэша")
	}
	user_id, _ := strconv.Atoi(val["user_id"])
	return user_id, val["login"], nil
}
