package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/CloudDetail/apo-sandbox/logging"
	"github.com/CloudDetail/apo-sandbox/model"
	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedis(addr string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisClient{Client: rdb}
}

func (c *RedisClient) SetUser(user *model.User, expiration time.Duration) error {
	key := "user:" + user.ID
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.Client.Set(context.Background(), key, userData, expiration).Err()
}

func (c *RedisClient) GetUser(id string) (*model.User, error) {
	key := "user:" + id
	userData, err := c.Client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user model.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *RedisClient) SetUserIDs(userIDs []string, expiration time.Duration) error {
	userIDsJSON, err := json.Marshal(userIDs)
	if err != nil {
		return err
	}
	return c.Client.Set(context.Background(), "all_user_ids", userIDsJSON, expiration).Err()
}

func (c *RedisClient) GetUserIDs() ([]string, error) {
	userIDsData, err := c.Client.Get(context.Background(), "all_user_ids").Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var userIDs []string
	err = json.Unmarshal([]byte(userIDsData), &userIDs)
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}

func (c *RedisClient) StartFault(delay int) error {
	cmd := c.Client.Do(context.Background(), "FAULT.START", delay)
	if cmd.Err() != nil {
		logging.Error("Failed to send fault command to Redis proxy: %v", cmd.Err())
		return cmd.Err()
	}
	logging.Info("Start redis latency %d.", delay)
	return nil
}

func (c *RedisClient) StopFault() error {
	cmd := c.Client.Do(context.Background(), "FAULT.STOP")
	if cmd.Err() != nil {
		logging.Error("Failed to send stop fault command to Redis proxy: %v", cmd.Err())
		return cmd.Err()
	}
	logging.Info("Redis latency stopped")
	return nil
}
