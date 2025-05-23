// Code generated by MockGen. DO NOT EDIT.
// Source: application/service/redis/client.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIRedisClient is a mock of IRedisClient interface.
type MockIRedisClient struct {
	ctrl     *gomock.Controller
	recorder *MockIRedisClientMockRecorder
}

// MockIRedisClientMockRecorder is the mock recorder for MockIRedisClient.
type MockIRedisClientMockRecorder struct {
	mock *MockIRedisClient
}

// NewMockIRedisClient creates a new mock instance.
func NewMockIRedisClient(ctrl *gomock.Controller) *MockIRedisClient {
	mock := &MockIRedisClient{ctrl: ctrl}
	mock.recorder = &MockIRedisClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIRedisClient) EXPECT() *MockIRedisClientMockRecorder {
	return m.recorder
}

// DeleteKey mocks base method.
func (m *MockIRedisClient) DeleteKey(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteKey", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteKey indicates an expected call of DeleteKey.
func (mr *MockIRedisClientMockRecorder) DeleteKey(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteKey", reflect.TypeOf((*MockIRedisClient)(nil).DeleteKey), key)
}

// GetKeysByPattern mocks base method.
func (m *MockIRedisClient) GetKeysByPattern(pattern string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKeysByPattern", pattern)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKeysByPattern indicates an expected call of GetKeysByPattern.
func (mr *MockIRedisClientMockRecorder) GetKeysByPattern(pattern interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKeysByPattern", reflect.TypeOf((*MockIRedisClient)(nil).GetKeysByPattern), pattern)
}

// GetValue mocks base method.
func (m *MockIRedisClient) GetValue(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValue indicates an expected call of GetValue.
func (mr *MockIRedisClientMockRecorder) GetValue(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockIRedisClient)(nil).GetValue), key)
}

// IncrementCounter mocks base method.
func (m *MockIRedisClient) IncrementCounter(key string, value float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementCounter", key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementCounter indicates an expected call of IncrementCounter.
func (mr *MockIRedisClientMockRecorder) IncrementCounter(key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementCounter", reflect.TypeOf((*MockIRedisClient)(nil).IncrementCounter), key, value)
}
