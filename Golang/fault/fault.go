package fault

import (
	"fmt"
	"sync"
)

type ChaosFault interface {
	Start(params map[string]interface{}) error
	Stop() error
	IsActive() bool
	Name() string
}

type Manager struct {
	mu     sync.Mutex
	faults map[string]ChaosFault
}

func NewManager() *Manager {
	return &Manager{
		faults: make(map[string]ChaosFault),
	}
}

func (m *Manager) Register(f ChaosFault) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.faults[f.Name()] = f
}

func (m *Manager) StartFault(chaosType string, params map[string]interface{}) error {
	m.mu.Lock()
	f, ok := m.faults[chaosType]
	m.mu.Unlock()
	if !ok {
		return fmt.Errorf("fault %s not found", chaosType)
	}
	return f.Start(params)
}

func (m *Manager) StopFault(chaosType string) error {
	m.mu.Lock()
	f, ok := m.faults[chaosType]
	m.mu.Unlock()
	if !ok {
		return fmt.Errorf("fault %s not found", chaosType)
	}
	return f.Stop()
}

func (m *Manager) StopAllFaults() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, f := range m.faults {
		f.Stop()
	}
}

func (m *Manager) Status() map[string]interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()

	status := make(map[string]interface{})
	for name, f := range m.faults {
		status[name] = map[string]interface{}{
			"active": f.IsActive(),
			"name":   f.Name(),
		}
	}
	return status
}

func (m *Manager) ListActive() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	var active []string
	for name, f := range m.faults {
		if f.IsActive() {
			active = append(active, name)
		}
	}
	return active
}
