package service

import (
	"encoding/json"
	"fmt"

	"github.com/CloudDetail/apo-sandbox/fault"
	"github.com/CloudDetail/apo-sandbox/logging"
	"github.com/CloudDetail/apo-sandbox/model"
	"github.com/CloudDetail/apo-sandbox/storage"
)

type BusinessService struct {
	Store        *storage.Store
	FaultManager *fault.Manager
}

func NewBusinessService(store *storage.Store, faultManager *fault.Manager) *BusinessService {
	return &BusinessService{
		Store:        store,
		FaultManager: faultManager,
	}
}

func (s *BusinessService) GetUsersCached(chaosType string, duration int) (string, error) {
	if len(chaosType) > 0 {
		params := map[string]interface{}{}
		if duration > 0 {
			params["duration"] = duration
		}
		err := s.FaultManager.StartFault(chaosType, params)
		if err != nil {
			logging.Error("Start fault failed: %v", err)
		}
	} else {
		s.FaultManager.StopAllFaults()
	}

	var users []model.User
	var err error

	if s.Store.Redis != nil {
		users, err = s.Store.QueryUsersCached()
		if err != nil {
			return "", err
		}
	}

	users, err = s.Store.QueryUsersWithDBBackup()
	if err != nil {
		return "", err
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		return "", fmt.Errorf("failed to marshal users: %w", err)
	}
	return string(usersJSON), nil
}
