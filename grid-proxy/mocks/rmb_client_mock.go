// Code generated by MockGen. DO NOT EDIT.
// Source: rmb-sdk-go/interface.go
//
// Generated by this command:
//
//	mockgen -source=rmb-sdk-go/interface.go -destination=grid-proxy/mocks/rmb_client_mock.go -typed
//

// Package mock_rmb is a generated GoMock package.
package mock_rmb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)


// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Call mocks base method.
func (m *MockClient) Call(ctx context.Context, twin uint32, fn string, data, result any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Call", ctx, twin, fn, data, result)
	ret0, _ := ret[0].(error)
	return ret0
}

// Call indicates an expected call of Call.
func (mr *MockClientMockRecorder) Call(ctx, twin, fn, data, result any) *MockClientCallCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Call", reflect.TypeOf((*MockClient)(nil).Call), ctx, twin, fn, data, result)
	return &MockClientCallCall{Call: call}
}

// MockClientCallCall wrap *gomock.Call
type MockClientCallCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockClientCallCall) Return(arg0 error) *MockClientCallCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockClientCallCall) Do(f func(context.Context, uint32, string, any, any) error) *MockClientCallCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockClientCallCall) DoAndReturn(f func(context.Context, uint32, string, any, any) error) *MockClientCallCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
