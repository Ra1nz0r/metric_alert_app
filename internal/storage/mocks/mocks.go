// Code generated by MockGen. DO NOT EDIT.
// Source: internal/storage/metrics.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetricService is a mock of MetricService interface.
type MockMetricService struct {
	ctrl     *gomock.Controller
	recorder *MockMetricServiceMockRecorder
}

// MockMetricServiceMockRecorder is the mock recorder for MockMetricService.
type MockMetricServiceMockRecorder struct {
	mock *MockMetricService
}

// NewMockMetricService creates a new mock instance.
func NewMockMetricService(ctrl *gomock.Controller) *MockMetricService {
	mock := &MockMetricService{ctrl: ctrl}
	mock.recorder = &MockMetricServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricService) EXPECT() *MockMetricServiceMockRecorder {
	return m.recorder
}

// MakeStorageCopy mocks base method.
func (m *MockMetricService) MakeStorageCopy() (*map[string]float64, *map[string]int64) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeStorageCopy")
	ret0, _ := ret[0].(*map[string]float64)
	ret1, _ := ret[1].(*map[string]int64)
	return ret0, ret1
}

// MakeStorageCopy indicates an expected call of MakeStorageCopy.
func (mr *MockMetricServiceMockRecorder) MakeStorageCopy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeStorageCopy", reflect.TypeOf((*MockMetricService)(nil).MakeStorageCopy))
}

// UpdateCounter mocks base method.
func (m *MockMetricService) UpdateCounter(name string, value int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateCounter", name, value)
}

// UpdateCounter indicates an expected call of UpdateCounter.
func (mr *MockMetricServiceMockRecorder) UpdateCounter(name, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCounter", reflect.TypeOf((*MockMetricService)(nil).UpdateCounter), name, value)
}

// UpdateGauge mocks base method.
func (m *MockMetricService) UpdateGauge(name string, value float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateGauge", name, value)
}

// UpdateGauge indicates an expected call of UpdateGauge.
func (mr *MockMetricServiceMockRecorder) UpdateGauge(name, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGauge", reflect.TypeOf((*MockMetricService)(nil).UpdateGauge), name, value)
}
