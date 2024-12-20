// Code generated by MockGen. DO NOT EDIT.
// Source: usecase/interfaces.go
//
// Generated by this command:
//
//	mockgen -destination=mocks/usecase.go -source=usecase/interfaces.go -package=mocks UsecaseInterface
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	generated "github.com/SawitProRecruitment/UserService/generated"
	model "github.com/SawitProRecruitment/UserService/model"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthUsecaseInterface is a mock of AuthUsecaseInterface interface.
type MockAuthUsecaseInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUsecaseInterfaceMockRecorder
	isgomock struct{}
}

// MockAuthUsecaseInterfaceMockRecorder is the mock recorder for MockAuthUsecaseInterface.
type MockAuthUsecaseInterfaceMockRecorder struct {
	mock *MockAuthUsecaseInterface
}

// NewMockAuthUsecaseInterface creates a new mock instance.
func NewMockAuthUsecaseInterface(ctrl *gomock.Controller) *MockAuthUsecaseInterface {
	mock := &MockAuthUsecaseInterface{ctrl: ctrl}
	mock.recorder = &MockAuthUsecaseInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUsecaseInterface) EXPECT() *MockAuthUsecaseInterfaceMockRecorder {
	return m.recorder
}

// LoginUser mocks base method.
func (m *MockAuthUsecaseInterface) LoginUser(ctx context.Context, payload generated.AuthLoginJSONRequestBody) (model.User, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", ctx, payload)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockAuthUsecaseInterfaceMockRecorder) LoginUser(ctx, payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockAuthUsecaseInterface)(nil).LoginUser), ctx, payload)
}

// MockUserUsecaseInterface is a mock of UserUsecaseInterface interface.
type MockUserUsecaseInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseInterfaceMockRecorder
	isgomock struct{}
}

// MockUserUsecaseInterfaceMockRecorder is the mock recorder for MockUserUsecaseInterface.
type MockUserUsecaseInterfaceMockRecorder struct {
	mock *MockUserUsecaseInterface
}

// NewMockUserUsecaseInterface creates a new mock instance.
func NewMockUserUsecaseInterface(ctrl *gomock.Controller) *MockUserUsecaseInterface {
	mock := &MockUserUsecaseInterface{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUsecaseInterface) EXPECT() *MockUserUsecaseInterfaceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserUsecaseInterface) CreateUser(ctx context.Context, payload generated.RegisterUserJSONRequestBody) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, payload)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserUsecaseInterfaceMockRecorder) CreateUser(ctx, payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserUsecaseInterface)(nil).CreateUser), ctx, payload)
}

// GetUserProfile mocks base method.
func (m *MockUserUsecaseInterface) GetUserProfile(ctx context.Context, userId int64) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserProfile", ctx, userId)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserProfile indicates an expected call of GetUserProfile.
func (mr *MockUserUsecaseInterfaceMockRecorder) GetUserProfile(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserProfile", reflect.TypeOf((*MockUserUsecaseInterface)(nil).GetUserProfile), ctx, userId)
}

// UpdateUserProfile mocks base method.
func (m *MockUserUsecaseInterface) UpdateUserProfile(ctx context.Context, userId int64, payload generated.UpdateUserProfileJSONRequestBody) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserProfile", ctx, userId, payload)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUserProfile indicates an expected call of UpdateUserProfile.
func (mr *MockUserUsecaseInterfaceMockRecorder) UpdateUserProfile(ctx, userId, payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserProfile", reflect.TypeOf((*MockUserUsecaseInterface)(nil).UpdateUserProfile), ctx, userId, payload)
}
