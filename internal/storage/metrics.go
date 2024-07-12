package storage

import (
	"fmt"
	"sync"
)

type MetricService interface {
	AllMetricsFromStorage() map[string]any
	GetMap() (*map[string]float64, *map[string]int64)
	GetMetricVal(mType, mName string) (map[string]any, error)
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

func (m *MemStorage) AllMetricsFromStorage() map[string]any {
	res := make(map[string]any)

	m.mu.Lock()
	defer m.mu.Unlock()

	for k, v := range m.gauge {
		res[k] = v
	}

	for k, v := range m.counter {
		res[k] = v
	}

	return res
}

func (m *MemStorage) GetMap() (*map[string]float64, *map[string]int64) {
	return &m.gauge, &m.counter
}

func (m *MemStorage) GetMetricVal(mType, mName string) (map[string]any, error) {
	res := make(map[string]any)

	switch mType {
	case "gauge":
		gVal, ok := m.gauge[mName]
		if ok {
			res[mName] = gVal
			return res, nil
		}
	case "counter":
		cVal, ok := m.counter[mName]
		if ok {
			res[mName] = cVal
			return res, nil
		}
	default:
		return nil, fmt.Errorf("type not found")
	}
	return nil, fmt.Errorf("metric not found")
}

func (m *MemStorage) MakeStorageCopy() (*map[string]float64, *map[string]int64) {
	newStrg := New()

	for k, v := range m.gauge {
		m.mu.RLock()
		defer m.mu.RUnlock()

		newStrg.gauge[k] = v
	}

	for k, v := range m.counter {
		m.mu.RLock()
		defer m.mu.RUnlock()

		newStrg.counter[k] = v
	}

	return &newStrg.gauge, &newStrg.counter
}

func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.gauge[name] = value
}

func (m *MemStorage) UpdateCounter(name string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counter[name] += value
}
