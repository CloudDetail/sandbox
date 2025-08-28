package fault

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/CloudDetail/apo-sandbox/config"
)

type LatencyFault struct {
	mu     sync.Mutex
	active bool
	iface  string
	delay  int // ms
}

func NewLatencyFault(iface string) *LatencyFault {
	return &LatencyFault{iface: iface}
}

func (l *LatencyFault) Start(params map[string]interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.active {
		return fmt.Errorf("latency fault already active")
	}
	delayMs := config.LoadConfig().Faults.Latency.DefaultDelay
	if delay, ok := params["duration"]; ok {
		delayMs = delay.(int)
	}
	if delayMs < 1 {
		delayMs = 100
	}

	_ = l.clearTC()

	cmd := exec.Command("tc", "qdisc", "add", "dev", l.iface, "root", "netem", "delay", fmt.Sprintf("%dms", delayMs))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("add delay failed: %v, output: %s", err, string(output))
	}

	l.delay = delayMs
	l.active = true
	return nil
}

func (l *LatencyFault) Stop() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.active {
		return nil
	}
	if err := l.clearTC(); err != nil {
		return err
	}
	l.active = false
	return nil
}

func (l *LatencyFault) clearTC() error {
	cmd := exec.Command("tc", "qdisc", "del", "dev", l.iface, "root")
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "No such file or directory") &&
		!strings.Contains(string(output), "No qdisc") {
		return fmt.Errorf("clear tc failed: %v, output: %s", err, string(output))
	}
	return nil
}

func (l *LatencyFault) Name() string {
	return "latency"
}

func (f *LatencyFault) IsActive() bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.active
}
