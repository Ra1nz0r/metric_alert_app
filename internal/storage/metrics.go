package storage

import (
	"sync"
)

type MetricService interface {
	AllMetricsFromStorage() map[string]any
	GetMap() (*map[string]float64, *map[string]int64)
	MakeStorageCopy() (*map[string]float64, *map[string]int64)
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
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

func (m *MemStorage) AllMetricsFromStorage() map[string]any {
	zz := make(map[string]any)

	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range m.gauge {
		zz[k] = v
	}

	for k, v := range m.counter {
		zz[k] = v
	}

	return zz
}

func (m *MemStorage) MakeStorageCopy() (*map[string]float64, *map[string]int64) {
	newStrg := New()

	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.gauge {
		newStrg.gauge[k] = v
	}

	for k, v := range m.counter {
		newStrg.counter[k] = v
	}

	return &newStrg.gauge, &newStrg.counter
}

func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.gauge[name] = value

	return
}

func (m *MemStorage) UpdateCounter(name string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter[name] += value

	return
}
