package storage

import (
	"sync"
)

type MetricService interface {
	UpdateGauge(name string, value float64) (*map[string]float64, error)
	UpdateCounter(name string, value int64) (*map[string]int64, error)
	GetMap() (*map[string]float64, *map[string]int64)
	MakeStorageCopy() (*map[string]float64, *map[string]int64)
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mu      sync.RWMutex
}

func New() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (m *MemStorage) GetMap() (*map[string]float64, *map[string]int64) {
	return &m.gauge, &m.counter
}

func (m *MemStorage) MakeStorageCopy() (*map[string]float64, *map[string]int64) {
	newStrg := New()

	m.mu.RLock()
	for k, v := range m.gauge {
		newStrg.gauge[k] = v
	}
	m.mu.RUnlock()

	m.mu.RLock()
	for k, v := range m.counter {
		newStrg.counter[k] = v
	}
	m.mu.RUnlock()

	return &newStrg.gauge, &newStrg.counter
}

func (m *MemStorage) UpdateGauge(name string, value float64) (*map[string]float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.gauge[name] = value

	return &m.gauge, nil
}

func (m *MemStorage) UpdateCounter(name string, value int64) (*map[string]int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter[name] += value

	return &m.counter, nil
}
