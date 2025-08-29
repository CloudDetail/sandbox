package fault

import (
	"fmt"
	"sync"

	"github.com/CloudDetail/apo-sandbox/config"
	"github.com/CloudDetail/apo-sandbox/logging"
	"github.com/CloudDetail/apo-sandbox/storage"
)

type RedisLatencyFault struct {
	mu          sync.Mutex
	active      bool
	redisClient *storage.RedisClient
}

func NewRedisLatencyFault(redisClient *storage.RedisClient) *RedisLatencyFault {
	return &RedisLatencyFault{
		redisClient: redisClient,
	}
}

func (f *RedisLatencyFault) Start(params map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.active {
		return nil
	}

	defaultDelay := config.LoadConfig().Faults.Redis.DefaultDelay
	if delay, ok := params["duration"].(int); ok {
		defaultDelay = delay
	}

	err := f.redisClient.StartFault(defaultDelay)
	if err != nil {
		return fmt.Errorf("failed to start redis latency fault: %w", err)
	}

	f.active = true
	return nil
}

func (f *RedisLatencyFault) Stop() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.active {
		return nil
	}

	err := f.redisClient.StopFault()
	if err != nil {
		return fmt.Errorf("failed to stop redis latency fault: %w", err)
	}

	f.active = false
	logging.Info("Redis latency fault stopped.")
	return nil
}

func (f *RedisLatencyFault) IsActive() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.active
}

func (f *RedisLatencyFault) Name() string {
	return "redis_latency"
}
