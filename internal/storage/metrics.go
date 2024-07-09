package storage

import (
	"sync"
)

var mu sync.Mutex

type MetricService interface {
	UpdateGauge(name string, value float64) (*map[string]float64, error)
	UpdateCounter(name string, value int64) (*map[string]int64, error)
	GetMap() (*map[string]float64, *map[string]int64)
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
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

func (m *MemStorage) UpdateGauge(name string, value float64) (*map[string]float64, error) {
	mu.Lock()
	defer mu.Unlock()

	m.gauge[name] = value

	return &m.gauge, nil
}

func (m *MemStorage) UpdateCounter(name string, value int64) (*map[string]int64, error) {
	mu.Lock()
	defer mu.Unlock()

	m.counter[name] += value

	return &m.counter, nil
}
