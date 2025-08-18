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

// QueryUsersFromDB queries users directly from MySQL database
func (s *Store) QueryUsersFromDB() ([]model.User, error) {
	// Check if MySQL client is nil
	if s.MySQL == nil {
		logging.Error("MySQL client is nil")
		return nil, fmt.Errorf("mysql client is nil")
	}

	// Query users from MySQL
	rows, err := s.MySQL.Query("SELECT id, name, email FROM users")
	if err != nil {
		logging.Error("Failed to query users from MySQL: %v", err)
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			logging.Warn("Failed to scan user row: %v", err)
			continue
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		logging.Error("Error after reading rows: %v", err)
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	logging.Info("Successfully queried %d users from MySQL", len(users))
	return users, nil
}

// QueryUsersWithDBBackup queries users from database first, if none exist, creates new ones
func (s *Store) QueryUsersWithDBBackup() ([]model.User, error) {
	// First try to get users from database
	users, err := s.QueryUsersFromDB()
	if err != nil {
		logging.Error("Failed to query users from database: %v", err)
		return nil, err
	}

	// If no users in database, create new ones
	if len(users) == 0 {
		logging.Info("No users found in database. Creating 10 mock users.")
		for i := 0; i < 10; i++ {
			user := model.User{
				ID:    uuid.New().String(),
				Name:  fmt.Sprintf("Mock User DB %d", i+1),
				Email: fmt.Sprintf("mock_db%d@example.com", i+1),
			}
			users = append(users, user)
			// Insert user into database
			_, err := s.MySQL.Exec("INSERT INTO users (id, name, email) VALUES (?, ?, ?)",
				user.ID, user.Name, user.Email)
			if err != nil {
				logging.Warn("Failed to insert user %s into database: %v", user.ID, err)
			}
		}

		logging.Info("Created and stored 10 mock users in database")
	}

	return users, nil
}

func (s *Store) QueryUsersCached() ([]model.User, error) {
	// If Redis client is nil, simulate HTTP operation
	if s.Redis == nil {
		logging.Info("Redis client is nil. Simulating HTTP operation to fetch users.")
		// Simulate a network delay for HTTP operation
		time.Sleep(10 * time.Millisecond)
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
