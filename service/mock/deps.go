// Code generated by MockGen. DO NOT EDIT.
// Source: service/service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	models "github.com/akhripko/dummy/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockStorage is a mock of Storage interface
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Check mocks base method
func (m *MockStorage) Check() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check")
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check
func (mr *MockStorageMockRecorder) Check() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockStorage)(nil).Check))
}

// Hello mocks base method
func (m *MockStorage) Hello(name string) (*models.HelloMessage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hello", name)
	ret0, _ := ret[0].(*models.HelloMessage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hello indicates an expected call of Hello
func (mr *MockStorageMockRecorder) Hello(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hello", reflect.TypeOf((*MockStorage)(nil).Hello), name)
}

// MockCache is a mock of Cache interface
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// Check mocks base method
func (m *MockCache) Check() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check")
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check
func (mr *MockCacheMockRecorder) Check() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockCache)(nil).Check))
}

// Read mocks base method
func (m *MockCache) Read(name string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", name)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read
func (mr *MockCacheMockRecorder) Read(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockCache)(nil).Read), name)
}

// WriteTTL mocks base method
func (m *MockCache) WriteTTL(name, msg string, ttl int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteTTL", name, msg, ttl)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteTTL indicates an expected call of WriteTTL
func (mr *MockCacheMockRecorder) WriteTTL(name, msg, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteTTL", reflect.TypeOf((*MockCache)(nil).WriteTTL), name, msg, ttl)
}
