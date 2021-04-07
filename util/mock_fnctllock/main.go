// Code generated by MockGen. DO NOT EDIT.
// Source: main.go

// Package mock_fcntllock is a generated GoMock package.
package mock_fcntllock

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	fcntllock "opensvc.com/opensvc/util/fcntllock"
)

// MockReadWriteSeekCloser is a mock of ReadWriteSeekCloser interface.
type MockReadWriteSeekCloser struct {
	ctrl     *gomock.Controller
	recorder *MockReadWriteSeekCloserMockRecorder
}

// MockReadWriteSeekCloserMockRecorder is the mock recorder for MockReadWriteSeekCloser.
type MockReadWriteSeekCloserMockRecorder struct {
	mock *MockReadWriteSeekCloser
}

// NewMockReadWriteSeekCloser creates a new mock instance.
func NewMockReadWriteSeekCloser(ctrl *gomock.Controller) *MockReadWriteSeekCloser {
	mock := &MockReadWriteSeekCloser{ctrl: ctrl}
	mock.recorder = &MockReadWriteSeekCloserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReadWriteSeekCloser) EXPECT() *MockReadWriteSeekCloserMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockReadWriteSeekCloser) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockReadWriteSeekCloserMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockReadWriteSeekCloser)(nil).Close))
}

// Read mocks base method.
func (m *MockReadWriteSeekCloser) Read(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockReadWriteSeekCloserMockRecorder) Read(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockReadWriteSeekCloser)(nil).Read), p)
}

// Seek mocks base method.
func (m *MockReadWriteSeekCloser) Seek(offset int64, whence int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seek", offset, whence)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Seek indicates an expected call of Seek.
func (mr *MockReadWriteSeekCloserMockRecorder) Seek(offset, whence interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seek", reflect.TypeOf((*MockReadWriteSeekCloser)(nil).Seek), offset, whence)
}

// Write mocks base method.
func (m *MockReadWriteSeekCloser) Write(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockReadWriteSeekCloserMockRecorder) Write(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockReadWriteSeekCloser)(nil).Write), p)
}

// MockLocker is a mock of Locker interface.
type MockLocker struct {
	ctrl     *gomock.Controller
	recorder *MockLockerMockRecorder
}

// MockLockerMockRecorder is the mock recorder for MockLocker.
type MockLockerMockRecorder struct {
	mock *MockLocker
}

// NewMockLocker creates a new mock instance.
func NewMockLocker(ctrl *gomock.Controller) *MockLocker {
	mock := &MockLocker{ctrl: ctrl}
	mock.recorder = &MockLockerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocker) EXPECT() *MockLockerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockLocker) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockLockerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockLocker)(nil).Close))
}

// LockContext mocks base method.
func (m *MockLocker) LockContext(arg0 context.Context, arg1 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LockContext", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LockContext indicates an expected call of LockContext.
func (mr *MockLockerMockRecorder) LockContext(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockContext", reflect.TypeOf((*MockLocker)(nil).LockContext), arg0, arg1)
}

// Read mocks base method.
func (m *MockLocker) Read(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockLockerMockRecorder) Read(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockLocker)(nil).Read), p)
}

// Seek mocks base method.
func (m *MockLocker) Seek(offset int64, whence int) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Seek", offset, whence)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Seek indicates an expected call of Seek.
func (mr *MockLockerMockRecorder) Seek(offset, whence interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Seek", reflect.TypeOf((*MockLocker)(nil).Seek), offset, whence)
}

// UnLock mocks base method.
func (m *MockLocker) UnLock() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnLock")
	ret0, _ := ret[0].(error)
	return ret0
}

// UnLock indicates an expected call of UnLock.
func (mr *MockLockerMockRecorder) UnLock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnLock", reflect.TypeOf((*MockLocker)(nil).UnLock))
}

// Write mocks base method.
func (m *MockLocker) Write(p []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Write", p)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Write indicates an expected call of Write.
func (mr *MockLockerMockRecorder) Write(p interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Write", reflect.TypeOf((*MockLocker)(nil).Write), p)
}

// MockLockProvider is a mock of LockProvider interface.
type MockLockProvider struct {
	ctrl     *gomock.Controller
	recorder *MockLockProviderMockRecorder
}

// MockLockProviderMockRecorder is the mock recorder for MockLockProvider.
type MockLockProviderMockRecorder struct {
	mock *MockLockProvider
}

// NewMockLockProvider creates a new mock instance.
func NewMockLockProvider(ctrl *gomock.Controller) *MockLockProvider {
	mock := &MockLockProvider{ctrl: ctrl}
	mock.recorder = &MockLockProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLockProvider) EXPECT() *MockLockProviderMockRecorder {
	return m.recorder
}

// New mocks base method.
func (m *MockLockProvider) New(arg0 string) fcntllock.Locker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", arg0)
	ret0, _ := ret[0].(fcntllock.Locker)
	return ret0
}

// New indicates an expected call of New.
func (mr *MockLockProviderMockRecorder) New(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockLockProvider)(nil).New), arg0)
}
