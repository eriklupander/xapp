// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/imageloader/imageloader.go

// Package mock_imageloader is a generated GoMock package.
package mock_imageloader

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockImageLoader is a mock of ImageLoader interface
type MockImageLoader struct {
	ctrl     *gomock.Controller
	recorder *MockImageLoaderMockRecorder
}

// MockImageLoaderMockRecorder is the mock recorder for MockImageLoader
type MockImageLoaderMockRecorder struct {
	mock *MockImageLoader
}

// NewMockImageLoader creates a new mock instance
func NewMockImageLoader(ctrl *gomock.Controller) *MockImageLoader {
	mock := &MockImageLoader{ctrl: ctrl}
	mock.recorder = &MockImageLoaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockImageLoader) EXPECT() *MockImageLoaderMockRecorder {
	return m.recorder
}

// Load mocks base method
func (m *MockImageLoader) Load(url string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", url)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load
func (mr *MockImageLoaderMockRecorder) Load(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockImageLoader)(nil).Load), url)
}