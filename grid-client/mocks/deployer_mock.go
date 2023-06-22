// Code generated by MockGen. DO NOT EDIT.
// Source: ./deployer/deployer.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	gridtypes "github.com/threefoldtech/zos/pkg/gridtypes"
)

// MockMockDeployer is a mock of MockDeployer interface.
type MockMockDeployer struct {
	ctrl     *gomock.Controller
	recorder *MockMockDeployerMockRecorder
}

// MockMockDeployerMockRecorder is the mock recorder for MockMockDeployer.
type MockMockDeployerMockRecorder struct {
	mock *MockMockDeployer
}

// NewMockMockDeployer creates a new mock instance.
func NewMockMockDeployer(ctrl *gomock.Controller) *MockMockDeployer {
	mock := &MockMockDeployer{ctrl: ctrl}
	mock.recorder = &MockMockDeployerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMockDeployer) EXPECT() *MockMockDeployerMockRecorder {
	return m.recorder
}

// BatchDeploy mocks base method.
func (m *MockMockDeployer) BatchDeploy(ctx context.Context, deployments map[uint32][]gridtypes.Deployment, deploymentsSolutionProvider map[uint32][]*uint64) (map[uint32][]gridtypes.Deployment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchDeploy", ctx, deployments, deploymentsSolutionProvider)
	ret0, _ := ret[0].(map[uint32][]gridtypes.Deployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BatchDeploy indicates an expected call of BatchDeploy.
func (mr *MockMockDeployerMockRecorder) BatchDeploy(ctx, deployments, deploymentsSolutionProvider interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchDeploy", reflect.TypeOf((*MockMockDeployer)(nil).BatchDeploy), ctx, deployments, deploymentsSolutionProvider)
}

// Cancel mocks base method.
func (m *MockMockDeployer) Cancel(ctx context.Context, contractID []uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Cancel", ctx, contractID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Cancel indicates an expected call of Cancel.
func (mr *MockMockDeployerMockRecorder) Cancel(ctx, contractID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Cancel", reflect.TypeOf((*MockMockDeployer)(nil).Cancel), ctx, contractID)
}

// Deploy mocks base method.
func (m *MockMockDeployer) Deploy(ctx context.Context, oldDeploymentIDs map[uint32]uint64, newDeployments map[uint32]gridtypes.Deployment, newDeploymentSolutionProvider map[uint32]*uint64) (map[uint32]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deploy", ctx, oldDeploymentIDs, newDeployments, newDeploymentSolutionProvider)
	ret0, _ := ret[0].(map[uint32]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Deploy indicates an expected call of Deploy.
func (mr *MockMockDeployerMockRecorder) Deploy(ctx, oldDeploymentIDs, newDeployments, newDeploymentSolutionProvider interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deploy", reflect.TypeOf((*MockMockDeployer)(nil).Deploy), ctx, oldDeploymentIDs, newDeployments, newDeploymentSolutionProvider)
}

// GetDeployments mocks base method.
func (m *MockMockDeployer) GetDeployments(ctx context.Context, dls map[uint32]uint64) (map[uint32]gridtypes.Deployment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDeployments", ctx, dls)
	ret0, _ := ret[0].(map[uint32]gridtypes.Deployment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDeployments indicates an expected call of GetDeployments.
func (mr *MockMockDeployerMockRecorder) GetDeployments(ctx, dls interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDeployments", reflect.TypeOf((*MockMockDeployer)(nil).GetDeployments), ctx, dls)
}
