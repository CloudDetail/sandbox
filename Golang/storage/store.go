package storage

import (
	"fmt"
	"time"

	"github.com/CloudDetail/apo-sandbox/logging"
	"github.com/CloudDetail/apo-sandbox/model"
	"github.com/google/uuid"
)

type Store struct {
	MySQL *MySQLClient
	Redis *RedisClient
}

func (s *Store) QueryUsersCached() ([]model.User, error) {
	// If Redis client is nil, simulate HTTP operation
	if s.Redis == nil {
		logging.Info("Redis client is nil. Simulating HTTP operation to fetch users.")
		// Simulate a network delay for HTTP operation
		time.Sleep(100 * time.Millisecond)
		var users []model.User
		for i := 0; i < 10; i++ {
			user := model.User{
				ID:    uuid.New().String(),
				Name:  fmt.Sprintf("Mock User HTTP %d", i+1),
				Email: fmt.Sprintf("mock_http%d@example.com", i+1),
			}
			users = append(users, user)
		}
		logging.Info("Successfully simulated HTTP operation for 10 users.")
		return users, nil
	}

	// Try to get user IDs from Redis cache first
	userIDs, err := s.Redis.GetUserIDs()
	if err == nil && len(userIDs) > 0 {
		var users []model.User
		for _, userID := range userIDs {
			user, err := s.Redis.GetUser(userID)
			if err != nil {
				logging.Warn("Failed to get user %s from Redis: %v", userID, err)
				continue
			}
			if user != nil {
				users = append(users, *user)
			}
		}
		if len(users) == len(userIDs) {
			logging.Info("%s", "All users retrieved from Redis cache by individual IDs.")
			return users, nil
		}
		logging.Warn("Incomplete users retrieved from Redis cache. Re-fetching from MySQL.")
	}

	// If not in Redis, mock 10 users and cache them
	var users []model.User
	var newserIDs []string
	for i := 0; i < 10; i++ {
		user := model.User{
			ID:    uuid.New().String(),
			Name:  fmt.Sprintf("Mock User %d", i+1),
			Email: fmt.Sprintf("mock%d@example.com", i+1),
		}
		users = append(users, user)
		newserIDs = append(newserIDs, user.ID)

		// Cache individual user
		err := s.Redis.SetUser(&user, 0)
		if err != nil {
			logging.Warn("Failed to cache user %s in Redis: %v", user.ID, err)
		}
	}

	// Cache all user IDs
	err = s.Redis.SetUserIDs(newserIDs, 0)
	if err != nil {
		logging.Warn("Failed to cache all user IDs in Redis: %v", err)
	}

	logging.Info("%s", "Mocked 10 users and cached in Redis (individual users and IDs).")
	return users, nil
}
