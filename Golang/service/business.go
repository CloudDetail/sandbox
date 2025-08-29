package service

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/CloudDetail/apo-sandbox/fault"
	"github.com/CloudDetail/apo-sandbox/storage"
	toxiproxy "github.com/Shopify/toxiproxy/v2/client"
)

type BusinessService struct {
	Store         *storage.Store
	FaultManager  *fault.Manager
	LatencyActive bool
	RedisActive   bool
	lmu           sync.Mutex
	rmu           sync.Mutex
}

func NewBusinessService(store *storage.Store) *BusinessService {
	return &BusinessService{
		Store: store,
	}
}

func (s *BusinessService) GetUsers1(mode string, duration int) (string, error) {
	if mode == "1" {
		s.lmu.Lock()
		defer s.lmu.Unlock()
		if !s.LatencyActive {
			_ = clearTC()
			// Injecting delay faults into the network adapter to simulate network latency
			cmd := exec.Command("tc", "qdisc", "add", "dev", "eth0", "root", "netem", "delay", fmt.Sprintf("%dms", duration))
			if _, err := cmd.CombinedOutput(); err != nil {
				return "", fmt.Errorf("type 1 failed")
			}

			s.LatencyActive = true
		}
	} else {
		err := s.stopFaults()
		if err != nil {
			return "", err
		}
	}
	// normal business
	s.Store.QueryUserFromRedis()
	users, err := s.Store.QueryUserFromMySQL()
	if err != nil {
		return "", err
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		return "", fmt.Errorf("failed to marshal users: %w", err)
	}
	return string(usersJSON), nil
}

func clearTC() error {
	cmd := exec.Command("tc", "qdisc", "del", "dev", "eth0", "root")
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "No such file or directory") &&
		!strings.Contains(string(output), "No qdisc") {
		return fmt.Errorf("clear tc failed: %v, output: %s", err, string(output))
	}
	return nil
}

func (s *BusinessService) GetUsers2(mode string, duration int) (string, error) {
	if mode == "1" {
		targetDuration := time.Duration(duration) * time.Millisecond

		start := time.Now()
		// Perform CPU-intensive operations to simulate CPU delays
		for time.Since(start) < targetDuration {
			_ = fibonacci(38)
		}
	} else {
		err := s.stopFaults()
		if err != nil {
			return "", err
		}
	}

	// noraml business
	s.Store.QueryUserFromRedis()
	users, err := s.Store.QueryUserFromMySQL()
	if err != nil {
		return "", err
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		return "", fmt.Errorf("failed to marshal users: %w", err)
	}
	return string(usersJSON), nil
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func (s *BusinessService) GetUsers3(mode string, duration int) (string, error) {
	if mode == "1" {
		s.rmu.Lock()
		defer s.rmu.Unlock()

		if !s.RedisActive {
			// Use https://github.com/Shopify/toxiproxy to simulate redis latency
			_, err := s.Store.Proxy.AddToxic(
				"redis_delay",
				"latency",
				"downstream",
				1.0,
				toxiproxy.Attributes{
					"latency": duration,
				},
			)
			if err != nil {
				return "", fmt.Errorf("type 3 failed")
			}
			s.RedisActive = true
		}
	} else {
		err := s.stopFaults()
		if err != nil {
			return "", err
		}
	}

	// noraml business
	s.Store.QueryUserFromRedis()
	users, err := s.Store.QueryUserFromMySQL()
	if err != nil {
		return "", err
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		return "", fmt.Errorf("failed to marshal users: %w", err)
	}
	return string(usersJSON), nil
}

func (s *BusinessService) stopFaults() error {
	s.lmu.Lock()
	defer s.lmu.Unlock()
	s.rmu.Lock()
	defer s.rmu.Unlock()

	if s.LatencyActive {
		if err := clearTC(); err != nil {
			return err
		}
		s.LatencyActive = false
	}

	if s.RedisActive {
		if err := s.Store.Proxy.RemoveToxic("redis_delay"); err != nil {
			return err
		}
		s.RedisActive = false
	}

	return nil
}
