package service

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/CloudDetail/apo-sandbox/fault"
	"github.com/CloudDetail/apo-sandbox/logging"
	"github.com/CloudDetail/apo-sandbox/storage"
	toxiproxy "github.com/Shopify/toxiproxy/v2/client"
)

type BusinessService struct {
	Store         *storage.Store
	FaultManager  *fault.Manager
	LatencyActive bool // Indicates if network latency fault is currently active
	RedisActive   bool // Indicates if Redis latency fault is currently active
	lmu           sync.Mutex
	rmu           sync.Mutex
}

func NewBusinessService(store *storage.Store) *BusinessService {
	return &BusinessService{
		Store: store,
	}
}

func (s *BusinessService) GetUsersWithLatency(mode string, duration int) (string, error) {
	if mode == "1" {
		s.lmu.Lock()
		defer s.lmu.Unlock()
		if !s.LatencyActive {
			// Clear any existing tc rules to ensure clean state
			_ = clearTC()

			// Inject network latency using Linux Traffic Control (tc)
			// Adds artificial delay to eth0 interface to simulate degraded network conditions
			// The delay parameter specifies milliseconds of additional latency per packet
			cmd := exec.Command("tc", "qdisc", "add", "dev", "eth0", "root", "netem", "delay", fmt.Sprintf("%dms", duration))
			if _, err := cmd.CombinedOutput(); err != nil {
				logging.Error("type 1 failed: %v", err)
				return "", fmt.Errorf("type 1 failed")
			}

			s.LatencyActive = true
		}
	} else {
		// Mode "0" or any other value stops all active faults
		err := s.stopFaults()
		if err != nil {
			return "", err
		}
	}

	// Continue with normal business operations regardless of fault state
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

// clearTC removes all traffic control rules from eth0 interface
// Used to reset network conditions to baseline after fault injection
func clearTC() error {
	cmd := exec.Command("tc", "qdisc", "del", "dev", "eth0", "root")
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "No such file or directory") &&
		!strings.Contains(string(output), "No qdisc") {
		return fmt.Errorf("clear tc failed: %v, output: %s", err, string(output))
	}
	return nil
}

func (s *BusinessService) GetUsersWithCPUBurn(mode string, duration int) (string, error) {
	if mode == "1" {
		targetDuration := time.Duration(duration) * time.Millisecond

		start := time.Now()
		// Simulate CPU stress using recursive Fibonacci calculation
		// This creates artificial CPU load to test system behavior under resource constraints
		// The calculation continues until the specified duration is reached
		for time.Since(start) < targetDuration {
			_ = fibonacci(38) // 38 provides sufficient computational load
		}
	} else {
		// Mode "0" or any other value stops all active faults
		err := s.stopFaults()
		if err != nil {
			return "", err
		}
	}

	// Continue with normal business operations regardless of fault state
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

// fibonacci provides a computationally expensive recursive function
// Used for generating CPU load in chaos testing scenarios
// Note: This is intentionally inefficient for stress testing purposes
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func (s *BusinessService) GetUsersWithRedisLatency(mode string, duration int) (string, error) {
	if mode == "1" {
		s.rmu.Lock()
		defer s.rmu.Unlock()

		if !s.RedisActive {
			// Use Toxiproxy to simulate Redis latency by adding toxic behavior
			// Toxiproxy creates a proxy between the application and Redis
			// The "latency" toxic introduces configurable delays to Redis responses
			_, err := s.Store.Proxy.AddToxic(
				"redis_delay",
				"latency",
				"downstream",
				1.0,
				toxiproxy.Attributes{
					"latency": duration, // Milliseconds of delay to inject
				},
			)
			if err != nil {
				logging.Error("type 3 failed: %v", err)
				return "", fmt.Errorf("type 3 failed")
			}
			s.RedisActive = true
		}
	} else {
		// Mode "0" or any other value stops all active faults
		err := s.stopFaults()
		if err != nil {
			return "", err
		}
	}

	// Continue with normal business operations regardless of fault state
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

// stopFaults terminates all active fault injection
// This provides a unified cleanup mechanism to restore normal system behavior
// Safely removes network latency rules, CPU stress, and Redis latency
func (s *BusinessService) stopFaults() error {
	s.lmu.Lock()
	defer s.lmu.Unlock()
	s.rmu.Lock()
	defer s.rmu.Unlock()

	// Remove network latency if active
	if s.LatencyActive {
		if err := clearTC(); err != nil {
			return err
		}
		s.LatencyActive = false
	}

	// Remove Redis latency if active
	if s.RedisActive {
		// Remove the latency toxic from Toxiproxy to restore normal Redis performance
		if err := s.Store.Proxy.RemoveToxic("redis_delay"); err != nil {
			return err
		}
		s.RedisActive = false
	}

	return nil
}
