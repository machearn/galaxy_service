// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/machearn/galaxy_service/mail (interfaces: EmailSender)

// Package mock_sender is a generated GoMock package.
package mock_sender

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEmailSender is a mock of EmailSender interface.
type MockEmailSender struct {
	ctrl     *gomock.Controller
	recorder *MockEmailSenderMockRecorder
}

// MockEmailSenderMockRecorder is the mock recorder for MockEmailSender.
type MockEmailSenderMockRecorder struct {
	mock *MockEmailSender
}

// NewMockEmailSender creates a new mock instance.
func NewMockEmailSender(ctrl *gomock.Controller) *MockEmailSender {
	mock := &MockEmailSender{ctrl: ctrl}
	mock.recorder = &MockEmailSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailSender) EXPECT() *MockEmailSenderMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockEmailSender) SendEmail(arg0 []string, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockEmailSenderMockRecorder) SendEmail(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockEmailSender)(nil).SendEmail), arg0, arg1, arg2)
}
