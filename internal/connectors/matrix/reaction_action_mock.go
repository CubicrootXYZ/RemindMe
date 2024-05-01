// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix (interfaces: ReactionAction)

// Package matrix is a generated GoMock package.
package matrix

import (
	reflect "reflect"

	gologger "github.com/CubicrootXYZ/gologger"
	database "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/database"
	mautrixcl "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/mautrixcl"
	messenger "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/connectors/matrix/messenger"
	database0 "github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/database"
	gomock "github.com/golang/mock/gomock"
)

// MockReactionAction is a mock of ReactionAction interface.
type MockReactionAction struct {
	ctrl     *gomock.Controller
	recorder *MockReactionActionMockRecorder
}

// MockReactionActionMockRecorder is the mock recorder for MockReactionAction.
type MockReactionActionMockRecorder struct {
	mock *MockReactionAction
}

// NewMockReactionAction creates a new mock instance.
func NewMockReactionAction(ctrl *gomock.Controller) *MockReactionAction {
	mock := &MockReactionAction{ctrl: ctrl}
	mock.recorder = &MockReactionActionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReactionAction) EXPECT() *MockReactionActionMockRecorder {
	return m.recorder
}

// Configure mocks base method.
func (m *MockReactionAction) Configure(arg0 gologger.Logger, arg1 mautrixcl.Client, arg2 messenger.Messenger, arg3 database.Service, arg4 database0.Service, arg5 *BridgeServices) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Configure", arg0, arg1, arg2, arg3, arg4, arg5)
}

// Configure indicates an expected call of Configure.
func (mr *MockReactionActionMockRecorder) Configure(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Configure", reflect.TypeOf((*MockReactionAction)(nil).Configure), arg0, arg1, arg2, arg3, arg4, arg5)
}

// HandleEvent mocks base method.
func (m *MockReactionAction) HandleEvent(arg0 *ReactionEvent, arg1 *database.MatrixMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleEvent", arg0, arg1)
}

// HandleEvent indicates an expected call of HandleEvent.
func (mr *MockReactionActionMockRecorder) HandleEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleEvent", reflect.TypeOf((*MockReactionAction)(nil).HandleEvent), arg0, arg1)
}

// Name mocks base method.
func (m *MockReactionAction) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockReactionActionMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockReactionAction)(nil).Name))
}

// Selector mocks base method.
func (m *MockReactionAction) Selector() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Selector")
	ret0, _ := ret[0].([]string)
	return ret0
}

// Selector indicates an expected call of Selector.
func (mr *MockReactionActionMockRecorder) Selector() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Selector", reflect.TypeOf((*MockReactionAction)(nil).Selector))
}