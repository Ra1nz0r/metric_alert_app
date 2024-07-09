package storage

import (
	"log"
	"strconv"
	"sync"
)

var mu sync.Mutex

type MetricService interface {
	UpdateGauge(name, value string) (*map[string]float64, error)
	UpdateCounter(name, value string) (*map[string]int64, error)
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

func (m *MemStorage) UpdateGauge(name, value string) (*map[string]float64, error) {
	mu.Lock()
	defer mu.Unlock()

	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Println("Error from strconv: ", err)
		return nil, err
	}
	m.gauge[name] = v

	return &m.gauge, nil
}

func (m *MemStorage) UpdateCounter(name, value string) (*map[string]int64, error) {
	mu.Lock()
	defer mu.Unlock()

	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Println("Error from strconv: ", err)
		return nil, err
	}
	m.counter[name] += v

	return &m.counter, nil
}
