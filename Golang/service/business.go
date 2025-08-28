package service

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/CloudDetail/apo-sandbox/config"
	"github.com/CloudDetail/apo-sandbox/fault"
	"github.com/CloudDetail/apo-sandbox/logging"
	"github.com/CloudDetail/apo-sandbox/model"
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

func NewBusinessService(store *storage.Store, faultManager *fault.Manager) *BusinessService {
	return &BusinessService{
		Store:        store,
		FaultManager: faultManager,
	}
}

func (s *BusinessService) GetUsers(chaosType string, duration int) (string, error) {
	switch chaosType {
	case "latency":
		// start network latency fault
		err := s.startLatenceFault(duration)
		if err != nil {
			return "", err
		}
	case "cpu":
		// start CPU latency fault
		err := s.startCPUFault(duration)
		if err != nil {
			return "", err
		}
	case "redis_latency":
		// start Redis latency fault
		err := s.startRedisFault(duration)
		if err != nil {
			return "", err
		}
	default:
		err := s.stopFaults()
		if err != nil {
			return "", err
		}
	}

	var users []model.User
	var err error

	s.Store.QueryUsersCached()
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

func (s *BusinessService) startCPUFault(duration int) error {
	faultConfig := config.LoadConfig().Faults.CPU
	targetDuration := time.Duration(faultConfig.DefaultDuration) * time.Millisecond
	if duration > 0 {
		targetDuration = time.Duration(duration) * time.Millisecond
	}

	start := time.Now()
	for time.Since(start) < targetDuration {
		_ = fibonacci(38)
	}

	logging.Info("CPU fault started, consumed %s CPU time.", time.Since(start).String())
	return nil
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func (s *BusinessService) startLatenceFault(duration int) error {
	s.lmu.Lock()
	defer s.lmu.Unlock()

	if s.LatencyActive {
		logging.Info("latency fault already active")
		return nil
	}
	delayMs := config.LoadConfig().Faults.Latency.DefaultDelay
	if duration > 0 {
		delayMs = duration
	}
	if delayMs < 1 {
		delayMs = 100
	}

	_ = clearTC()

	cmd := exec.Command("tc", "qdisc", "add", "dev", "eth0", "root", "netem", "delay", fmt.Sprintf("%dms", delayMs))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("add delay failed: %v, output: %s", err, string(output))
	}

	s.LatencyActive = true
	logging.Info("Successfully add %dms delay on %s", delayMs, "eth0")
	return nil
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

func (s *BusinessService) startRedisFault(duration int) error {
	s.rmu.Lock()
	defer s.rmu.Unlock()

	if s.RedisActive {
		logging.Info("redis fault already active")
		return nil
	}
	faultConfig := config.LoadConfig().Faults.Redis
	targetDuration := time.Duration(faultConfig.DefaultDelay)
	if duration > 0 {
		targetDuration = time.Duration(duration)
	}

	// use https://github.com/Shopify/toxiproxy to sumulate redis latency
	s.Store.Proxy.AddToxic(
		"redis_delay",
		"latency",
		"downstream",
		1.0,
		toxiproxy.Attributes{
			"latency": targetDuration,
		},
	)
	logging.Info("Redis fault started, delay %dms.", targetDuration)
	return nil
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
