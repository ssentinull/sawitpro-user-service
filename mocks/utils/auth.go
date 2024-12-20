// Code generated by MockGen. DO NOT EDIT.
// Source: utils/auth.go
//
// Generated by this command:
//
//	mockgen -destination=mocks/utils/auth.go -source=utils/auth.go -package=mocks AuthInterface
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	model "github.com/SawitProRecruitment/UserService/model"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthInterface is a mock of AuthInterface interface.
type MockAuthInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAuthInterfaceMockRecorder
	isgomock struct{}
}

// MockAuthInterfaceMockRecorder is the mock recorder for MockAuthInterface.
type MockAuthInterfaceMockRecorder struct {
	mock *MockAuthInterface
}

// NewMockAuthInterface creates a new mock instance.
func NewMockAuthInterface(ctrl *gomock.Controller) *MockAuthInterface {
	mock := &MockAuthInterface{ctrl: ctrl}
	mock.recorder = &MockAuthInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthInterface) EXPECT() *MockAuthInterfaceMockRecorder {
	return m.recorder
}

// GenerateJWTToken mocks base method.
func (m *MockAuthInterface) GenerateJWTToken(user model.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateJWTToken", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateJWTToken indicates an expected call of GenerateJWTToken.
func (mr *MockAuthInterfaceMockRecorder) GenerateJWTToken(user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateJWTToken", reflect.TypeOf((*MockAuthInterface)(nil).GenerateJWTToken), user)
}

// GetUserId mocks base method.
func (m *MockAuthInterface) GetUserId(tokenStr string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserId", tokenStr)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserId indicates an expected call of GetUserId.
func (mr *MockAuthInterfaceMockRecorder) GetUserId(tokenStr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserId", reflect.TypeOf((*MockAuthInterface)(nil).GetUserId), tokenStr)
}

// ValidateJWTToken mocks base method.
func (m *MockAuthInterface) ValidateJWTToken(tokenStr string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateJWTToken", tokenStr)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateJWTToken indicates an expected call of ValidateJWTToken.
func (mr *MockAuthInterfaceMockRecorder) ValidateJWTToken(tokenStr any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateJWTToken", reflect.TypeOf((*MockAuthInterface)(nil).ValidateJWTToken), tokenStr)
}
