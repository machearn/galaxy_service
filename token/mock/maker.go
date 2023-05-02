// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/machearn/galaxy_service/token (interfaces: TokenMaker)

// Package mock_token is a generated GoMock package.
package mock_token

import (
	reflect "reflect"
	time "time"

	paseto "aidanwoods.dev/go-paseto"
	gomock "github.com/golang/mock/gomock"
)

// MockTokenMaker is a mock of TokenMaker interface.
type MockTokenMaker struct {
	ctrl     *gomock.Controller
	recorder *MockTokenMakerMockRecorder
}

// MockTokenMakerMockRecorder is the mock recorder for MockTokenMaker.
type MockTokenMakerMockRecorder struct {
	mock *MockTokenMaker
}

// NewMockTokenMaker creates a new mock instance.
func NewMockTokenMaker(ctrl *gomock.Controller) *MockTokenMaker {
	mock := &MockTokenMaker{ctrl: ctrl}
	mock.recorder = &MockTokenMakerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenMaker) EXPECT() *MockTokenMakerMockRecorder {
	return m.recorder
}

// CreateToken mocks base method.
func (m *MockTokenMaker) CreateToken(arg0 int32, arg1 time.Duration) (string, *paseto.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*paseto.Token)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockTokenMakerMockRecorder) CreateToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockTokenMaker)(nil).CreateToken), arg0, arg1)
}

// VerifyToken mocks base method.
func (m *MockTokenMaker) VerifyToken(arg0 string) (*paseto.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyToken", arg0)
	ret0, _ := ret[0].(*paseto.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyToken indicates an expected call of VerifyToken.
func (mr *MockTokenMakerMockRecorder) VerifyToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockTokenMaker)(nil).VerifyToken), arg0)
}