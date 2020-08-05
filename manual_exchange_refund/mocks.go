// Code generated by MockGen. DO NOT EDIT.
// Source: manual_exchange_refund/comparison_order_importer.go

// Package manual_exchange_refund is a generated GoMock package.
package manual_exchange_refund

import (
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"reflect"
	frule_module "github.com/ifrolikov/fork-frule-module"
)

// MockComparisonOrderImporterInterface is a mock of ComparisonOrderImporterInterface interface
type MockComparisonOrderImporterInterface struct {
	ctrl     *gomock.Controller
	recorder *MockComparisonOrderImporterInterfaceMockRecorder
}

// MockComparisonOrderImporterInterfaceMockRecorder is the mock recorder for MockComparisonOrderImporterInterface
type MockComparisonOrderImporterInterfaceMockRecorder struct {
	mock *MockComparisonOrderImporterInterface
}

// NewMockComparisonOrderImporterInterface creates a new mock instance
func NewMockComparisonOrderImporterInterface(ctrl *gomock.Controller) *MockComparisonOrderImporterInterface {
	mock := &MockComparisonOrderImporterInterface{ctrl: ctrl}
	mock.recorder = &MockComparisonOrderImporterInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockComparisonOrderImporterInterface) EXPECT() *MockComparisonOrderImporterInterfaceMockRecorder {
	return m.recorder
}

// getComparisonOrder mocks base method
func (m *MockComparisonOrderImporterInterface) getComparisonOrder(logger zerolog.Logger) (frule_module.ComparisonOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getComparisonOrder", logger)
	ret0, _ := ret[0].(frule_module.ComparisonOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getComparisonOrder indicates an expected call of getComparisonOrder
func (mr *MockComparisonOrderImporterInterfaceMockRecorder) getComparisonOrder(logger interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getComparisonOrder", reflect.TypeOf((*MockComparisonOrderImporterInterface)(nil).getComparisonOrder), logger)
}

// MockComparisonOrderUpdaterInterface is a mock of ComparisonOrderUpdaterInterface interface
type MockComparisonOrderUpdaterInterface struct {
	ctrl     *gomock.Controller
	recorder *MockComparisonOrderUpdaterInterfaceMockRecorder
}

// MockComparisonOrderUpdaterInterfaceMockRecorder is the mock recorder for MockComparisonOrderUpdaterInterface
type MockComparisonOrderUpdaterInterfaceMockRecorder struct {
	mock *MockComparisonOrderUpdaterInterface
}

// NewMockComparisonOrderUpdaterInterface creates a new mock instance
func NewMockComparisonOrderUpdaterInterface(ctrl *gomock.Controller) *MockComparisonOrderUpdaterInterface {
	mock := &MockComparisonOrderUpdaterInterface{ctrl: ctrl}
	mock.recorder = &MockComparisonOrderUpdaterInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockComparisonOrderUpdaterInterface) EXPECT() *MockComparisonOrderUpdaterInterfaceMockRecorder {
	return m.recorder
}

// update mocks base method
func (m *MockComparisonOrderUpdaterInterface) update(logger zerolog.Logger) (*comparisonOrderContainer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "update", logger)
	ret0, _ := ret[0].(*comparisonOrderContainer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// update indicates an expected call of update
func (mr *MockComparisonOrderUpdaterInterfaceMockRecorder) update(logger interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "update", reflect.TypeOf((*MockComparisonOrderUpdaterInterface)(nil).update), logger)
}
