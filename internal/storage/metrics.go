package storage

import (
	"fmt"
	"sync"
)

//go:generate mockgen -source=internal/storage/metrics.go -destination=internal/storage/mocks/mocks.go -package=mocks

// Интерфейс для взаимодействия с локальным хранилищем.
type MetricService interface {
	AllMetricsFromStorage() *map[string]any
	GetMap() (*map[string]float64, *map[string]int64)
	GetMetricVal(mType, mName string) (any, error)
	MakeStorageCopy() (*map[string]float64, *map[string]int64)
	UpdateGauge(name string, value float64)
	UpdateCounter(name string, value int64)
}

// Локальное хранилище.
type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mu      sync.RWMutex
}

// Создаем и инициализируем новое хранилище.
func New() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

// Создает и возвращает указатель на новое хранилище, заполняя её
// всеми хранящимися метриками в локальном хранилище.
func (ms *MemStorage) AllMetricsFromStorage() *map[string]any {
	res := make(map[string]any)

	ms.mu.Lock()
	defer ms.mu.Unlock()

	for k, v := range ms.gauge {
		res[k] = float64(v)
	}

	for k, v := range ms.counter {
		res[k] = int64(v)
	}

	return &res
}

func (ms *MemStorage) GetMap() (*map[string]float64, *map[string]int64) {
	return &ms.gauge, &ms.counter
}

// Получает и возвращает значение метрики, в зависимости от указанного
// типа и имени. Если тип или имя метрики не найдено, возращает ошибку.
func (ms *MemStorage) GetMetricVal(mType, mName string) (any, error) {
	switch mType {
	case "gauge":
		gVal, ok := ms.gauge[mName]
		if ok {
			return gVal, nil
		}
		return nil, fmt.Errorf("metric not found")
	case "counter":
		cVal, ok := ms.counter[mName]
		if ok {
			return cVal, nil
		}
		return nil, fmt.Errorf("metric not found")
	default:
		return nil, fmt.Errorf("type not found")
	}
}

// Создает и инициализирует новое хранилище для метрик, заполняет его данными
// из локального хранилища. Возвращает указатель на новое хранилище.
func (ms *MemStorage) MakeStorageCopy() (*map[string]float64, *map[string]int64) {
	newStrg := New()

	for k, v := range ms.gauge {
		ms.mu.RLock()
		defer ms.mu.RUnlock()

		newStrg.gauge[k] = v
	}

	for k, v := range ms.counter {
		ms.mu.RLock()
		defer ms.mu.RUnlock()

		newStrg.counter[k] = v
	}

	return &newStrg.gauge, &newStrg.counter
}

// Обновляет и заменяет значение новым для метрик Gauge.
func (ms *MemStorage) UpdateGauge(name string, value float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.gauge[name] = value
}

// Обновляет и инкрементирует значение метрики Counter.
func (ms *MemStorage) UpdateCounter(name string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.counter[name] += value
}
