package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"AbdelrahmanDwedar/blogo/internal/domain/entity"
	"github.com/redis/go-redis/v9"
)

// RedisCache implements the cache repository interface
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache creates a new Redis cache connection
func NewRedisCache() *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Printf("⚠️  Warning: Redis connection failed: %v\n", err)
		fmt.Println("   Continuing without caching...")
		return nil
	}

	fmt.Println("✅ Redis cache connected successfully")
	return &RedisCache{
		client: client,
		ctx:    ctx,
	}
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	if r == nil || r.client == nil {
		return nil
	}
	return r.client.Close()
}

// SetUser caches a user
func (r *RedisCache) SetUser(user *entity.User, expiration time.Duration) error {
	if r == nil || r.client == nil {
		return nil
	}

	key := fmt.Sprintf("user:%d", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}

	return r.client.Set(r.ctx, key, data, expiration).Err()
}

// GetUser retrieves a cached user
func (r *RedisCache) GetUser(userID int64) (*entity.User, error) {
	if r == nil || r.client == nil {
		return nil, fmt.Errorf("cache not available")
	}

	key := fmt.Sprintf("user:%d", userID)
	data, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("user not in cache")
	}
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}

	return &user, nil
}

// DeleteUser removes a user from cache
func (r *RedisCache) DeleteUser(userID int64) error {
	if r == nil || r.client == nil {
		return nil
	}

	key := fmt.Sprintf("user:%d", userID)
	return r.client.Del(r.ctx, key).Err()
}

// SetBlog caches a blog
func (r *RedisCache) SetBlog(blog *entity.Blog, expiration time.Duration) error {
	if r == nil || r.client == nil {
		return nil
	}

	key := fmt.Sprintf("blog:%d", blog.ID)
	data, err := json.Marshal(blog)
	if err != nil {
		return fmt.Errorf("marshal blog: %w", err)
	}

	return r.client.Set(r.ctx, key, data, expiration).Err()
}

// GetBlog retrieves a cached blog
func (r *RedisCache) GetBlog(blogID int64) (*entity.Blog, error) {
	if r == nil || r.client == nil {
		return nil, fmt.Errorf("cache not available")
	}

	key := fmt.Sprintf("blog:%d", blogID)
	data, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("blog not in cache")
	}
	if err != nil {
		return nil, err
	}

	var blog entity.Blog
	if err := json.Unmarshal([]byte(data), &blog); err != nil {
		return nil, fmt.Errorf("unmarshal blog: %w", err)
	}

	return &blog, nil
}

// DeleteBlog removes a blog from cache
func (r *RedisCache) DeleteBlog(blogID int64) error {
	if r == nil || r.client == nil {
		return nil
	}

	key := fmt.Sprintf("blog:%d", blogID)
	return r.client.Del(r.ctx, key).Err()
}

// DeletePattern deletes all keys matching a pattern
func (r *RedisCache) DeletePattern(pattern string) error {
	if r == nil || r.client == nil {
		return nil
	}

	iter := r.client.Scan(r.ctx, 0, pattern, 0).Iterator()
	for iter.Next(r.ctx) {
		if err := r.client.Del(r.ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}


