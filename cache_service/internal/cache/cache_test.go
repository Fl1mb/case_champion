package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*Cache, func()) {
	mr, err := miniredis.Run()
	require.NoError(t, err)

	options := &redis.Options{
		Addr: mr.Addr(),
	}

	cache, err := NewCache(options)
	require.NoError(t, err)

	return cache, func() {
		mr.Close()
		cache.Close()
	}
}

func TestNewCache(t *testing.T) {
	t.Run("successful connection", func(t *testing.T) {
		mr, err := miniredis.Run()
		require.NoError(t, err)
		defer mr.Close()

		options := &redis.Options{
			Addr: mr.Addr(),
		}

		cache, err := NewCache(options)
		require.NoError(t, err)
		assert.NotNil(t, cache.rdb)
		cache.Close()
	})

	t.Run("failed connection", func(t *testing.T) {
		options := &redis.Options{
			Addr: "invalid-address:6379",
		}

		cache, err := NewCache(options)
		assert.Error(t, err)
		assert.Nil(t, cache)
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("healthy", func(t *testing.T) {
		cache, cleanup := setupTestRedis(t)
		defer cleanup()

		err := cache.HealthCheck(context.Background())
		assert.NoError(t, err)
	})

	t.Run("not initialized", func(t *testing.T) {
		var cache *Cache
		err := cache.HealthCheck(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "Cache is not init", err.Error())
	})

	t.Run("unhealthy", func(t *testing.T) {
		cache, cleanup := setupTestRedis(t)
		cleanup() // закрываем соединение

		err := cache.HealthCheck(context.Background())
		assert.Error(t, err)
	})
}

func TestSetData(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("successful set data", func(t *testing.T) {
		jwt := "test-jwt-token"
		login := "testuser"
		userID := int32(123)

		err := cache.SetData(ctx, jwt, login, userID)
		assert.NoError(t, err)

		// Проверяем, что данные записались
		val, err := cache.rdb.HGetAll(ctx, jwt).Result()
		assert.NoError(t, err)
		assert.Equal(t, login, val["login"])
		assert.Equal(t, "123", val["user_id"])
	})

	t.Run("empty jwt", func(t *testing.T) {
		err := cache.SetData(ctx, "", "user", 123)
		assert.Error(t, err)
	})
}

func TestSetDeadTime(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("successful set TTL", func(t *testing.T) {
		jwt := "test-jwt-ttl"
		err := cache.SetData(ctx, jwt, "user", 123)
		require.NoError(t, err)

		ttl := time.Minute * 30
		err = cache.SetDeadTime(ctx, jwt, ttl)
		assert.NoError(t, err)

		// Проверяем TTL
		actualTTL := cache.rdb.TTL(ctx, jwt).Val()
		assert.True(t, actualTTL > ttl-time.Second*5 && actualTTL <= ttl)
	})

	t.Run("non-existent key", func(t *testing.T) {
		err := cache.SetDeadTime(ctx, "non-existent-key", time.Minute)
		assert.NoError(t, err) // Redis не возвращает ошибку для несуществующих ключей
	})
}

func TestDelete(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		jwt := "test-jwt-delete"
		err := cache.SetData(ctx, jwt, "user", 123)
		require.NoError(t, err)

		// Проверяем, что ключ существует
		exists := cache.Exists(ctx, jwt)
		assert.True(t, exists)

		// Удаляем
		err = cache.Delete(ctx, jwt)
		assert.NoError(t, err)

		// Проверяем, что ключ удален
		exists = cache.Exists(ctx, jwt)
		assert.False(t, exists)
	})

	t.Run("delete non-existent key", func(t *testing.T) {
		err := cache.Delete(ctx, "non-existent-key")
		assert.NoError(t, err)
	})
}

func TestExists(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("key exists", func(t *testing.T) {
		jwt := "test-jwt-exists"
		err := cache.SetData(ctx, jwt, "user", 123)
		require.NoError(t, err)

		exists := cache.Exists(ctx, jwt)
		assert.True(t, exists)
	})

	t.Run("key does not exist", func(t *testing.T) {
		exists := cache.Exists(ctx, "non-existent-key")
		assert.False(t, exists)
	})
}

func TestGetData(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("successful get data", func(t *testing.T) {
		jwt := "test-jwt-get"
		login := "testuser"
		userID := 456

		err := cache.SetData(ctx, jwt, login, int32(userID))
		require.NoError(t, err)

		gotUserID, gotLogin, err := cache.GetData(ctx, jwt)
		assert.NoError(t, err)
		assert.Equal(t, userID, gotUserID)
		assert.Equal(t, login, gotLogin)
	})

	t.Run("non-existent key", func(t *testing.T) {
		_, _, err := cache.GetData(ctx, "non-existent-key")
		assert.Error(t, err)
	})
}

func TestClose(t *testing.T) {
	cache, cleanup := setupTestRedis(t)
	defer cleanup()

	err := cache.Close()
	assert.NoError(t, err)

	// Проверяем, что соединение действительно закрыто
	err = cache.HealthCheck(context.Background())
	assert.Error(t, err)
}
