package fault

import (
	"time"

	"github.com/CloudDetail/apo-sandbox/config"
)

type cpuFault struct {
}

func NewCPUFault() ChaosFault {
	return &cpuFault{}
}

func (c *cpuFault) Name() string {
	return "cpu"
}

func (c *cpuFault) Start(params map[string]interface{}) error {
	faultConfig := config.LoadConfig().Faults.CPU
	targetDuration := time.Duration(faultConfig.DefaultDuration) * time.Millisecond
	if duration, ok := params["duration"]; ok {
		targetDuration = time.Millisecond * time.Duration(duration.(int))
	}

	start := time.Now()
	for time.Since(start) < targetDuration {
		_ = fibonacci(38)
	}

	return nil
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func (c *cpuFault) Stop() error {
	return nil
}

func (c *cpuFault) IsActive() bool {
	return false
}
