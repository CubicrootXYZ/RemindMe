// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database (interfaces: InputService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockInputService is a mock of InputService interface.
type MockInputService struct {
	ctrl     *gomock.Controller
	recorder *MockInputServiceMockRecorder
}

// MockInputServiceMockRecorder is the mock recorder for MockInputService.
type MockInputServiceMockRecorder struct {
	mock *MockInputService
}

// NewMockInputService creates a new mock instance.
func NewMockInputService(ctrl *gomock.Controller) *MockInputService {
	mock := &MockInputService{ctrl: ctrl}
	mock.recorder = &MockInputServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInputService) EXPECT() *MockInputServiceMockRecorder {
	return m.recorder
}

// InputRemoved mocks base method.
func (m *MockInputService) InputRemoved(arg0 string, arg1 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InputRemoved", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InputRemoved indicates an expected call of InputRemoved.
func (mr *MockInputServiceMockRecorder) InputRemoved(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InputRemoved", reflect.TypeOf((*MockInputService)(nil).InputRemoved), arg0, arg1)
}
