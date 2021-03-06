// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/state (interfaces: EndpointBinding)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	state "github.com/juju/juju/state"
	reflect "reflect"
)

// MockEndpointBinding is a mock of EndpointBinding interface
type MockEndpointBinding struct {
	ctrl     *gomock.Controller
	recorder *MockEndpointBindingMockRecorder
}

// MockEndpointBindingMockRecorder is the mock recorder for MockEndpointBinding
type MockEndpointBindingMockRecorder struct {
	mock *MockEndpointBinding
}

// NewMockEndpointBinding creates a new mock instance
func NewMockEndpointBinding(ctrl *gomock.Controller) *MockEndpointBinding {
	mock := &MockEndpointBinding{ctrl: ctrl}
	mock.recorder = &MockEndpointBindingMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEndpointBinding) EXPECT() *MockEndpointBindingMockRecorder {
	return m.recorder
}

// DefaultEndpointBindingSpace mocks base method
func (m *MockEndpointBinding) DefaultEndpointBindingSpace() (string, error) {
	ret := m.ctrl.Call(m, "DefaultEndpointBindingSpace")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DefaultEndpointBindingSpace indicates an expected call of DefaultEndpointBindingSpace
func (mr *MockEndpointBindingMockRecorder) DefaultEndpointBindingSpace() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultEndpointBindingSpace", reflect.TypeOf((*MockEndpointBinding)(nil).DefaultEndpointBindingSpace))
}

// Space mocks base method
func (m *MockEndpointBinding) Space(arg0 string) (*state.Space, error) {
	ret := m.ctrl.Call(m, "Space", arg0)
	ret0, _ := ret[0].(*state.Space)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Space indicates an expected call of Space
func (mr *MockEndpointBindingMockRecorder) Space(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Space", reflect.TypeOf((*MockEndpointBinding)(nil).Space), arg0)
}

// SpaceIDsByName mocks base method
func (m *MockEndpointBinding) SpaceIDsByName() (map[string]string, error) {
	ret := m.ctrl.Call(m, "SpaceIDsByName")
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SpaceIDsByName indicates an expected call of SpaceIDsByName
func (mr *MockEndpointBindingMockRecorder) SpaceIDsByName() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpaceIDsByName", reflect.TypeOf((*MockEndpointBinding)(nil).SpaceIDsByName))
}

// SpaceNamesByID mocks base method
func (m *MockEndpointBinding) SpaceNamesByID() (map[string]string, error) {
	ret := m.ctrl.Call(m, "SpaceNamesByID")
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SpaceNamesByID indicates an expected call of SpaceNamesByID
func (mr *MockEndpointBindingMockRecorder) SpaceNamesByID() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpaceNamesByID", reflect.TypeOf((*MockEndpointBinding)(nil).SpaceNamesByID))
}
