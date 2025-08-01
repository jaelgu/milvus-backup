// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package taskmgr

import (
	"time"

	mock "github.com/stretchr/testify/mock"

	"github.com/zilliztech/milvus-backup/core/proto/backuppb"
)

// NewMockRestoreCollTaskView creates a new instance of MockRestoreCollTaskView. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRestoreCollTaskView(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRestoreCollTaskView {
	mock := &MockRestoreCollTaskView{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockRestoreCollTaskView is an autogenerated mock type for the RestoreCollTaskView type
type MockRestoreCollTaskView struct {
	mock.Mock
}

type MockRestoreCollTaskView_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRestoreCollTaskView) EXPECT() *MockRestoreCollTaskView_Expecter {
	return &MockRestoreCollTaskView_Expecter{mock: &_m.Mock}
}

// EndTime provides a mock function for the type MockRestoreCollTaskView
func (_mock *MockRestoreCollTaskView) EndTime() time.Time {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for EndTime")
	}

	var r0 time.Time
	if returnFunc, ok := ret.Get(0).(func() time.Time); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(time.Time)
	}
	return r0
}

// MockRestoreCollTaskView_EndTime_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EndTime'
type MockRestoreCollTaskView_EndTime_Call struct {
	*mock.Call
}

// EndTime is a helper method to define mock.On call
func (_e *MockRestoreCollTaskView_Expecter) EndTime() *MockRestoreCollTaskView_EndTime_Call {
	return &MockRestoreCollTaskView_EndTime_Call{Call: _e.mock.On("EndTime")}
}

func (_c *MockRestoreCollTaskView_EndTime_Call) Run(run func()) *MockRestoreCollTaskView_EndTime_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRestoreCollTaskView_EndTime_Call) Return(time1 time.Time) *MockRestoreCollTaskView_EndTime_Call {
	_c.Call.Return(time1)
	return _c
}

func (_c *MockRestoreCollTaskView_EndTime_Call) RunAndReturn(run func() time.Time) *MockRestoreCollTaskView_EndTime_Call {
	_c.Call.Return(run)
	return _c
}

// ErrorMessage provides a mock function for the type MockRestoreCollTaskView
func (_mock *MockRestoreCollTaskView) ErrorMessage() string {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for ErrorMessage")
	}

	var r0 string
	if returnFunc, ok := ret.Get(0).(func() string); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(string)
	}
	return r0
}

// MockRestoreCollTaskView_ErrorMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ErrorMessage'
type MockRestoreCollTaskView_ErrorMessage_Call struct {
	*mock.Call
}

// ErrorMessage is a helper method to define mock.On call
func (_e *MockRestoreCollTaskView_Expecter) ErrorMessage() *MockRestoreCollTaskView_ErrorMessage_Call {
	return &MockRestoreCollTaskView_ErrorMessage_Call{Call: _e.mock.On("ErrorMessage")}
}

func (_c *MockRestoreCollTaskView_ErrorMessage_Call) Run(run func()) *MockRestoreCollTaskView_ErrorMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRestoreCollTaskView_ErrorMessage_Call) Return(s string) *MockRestoreCollTaskView_ErrorMessage_Call {
	_c.Call.Return(s)
	return _c
}

func (_c *MockRestoreCollTaskView_ErrorMessage_Call) RunAndReturn(run func() string) *MockRestoreCollTaskView_ErrorMessage_Call {
	_c.Call.Return(run)
	return _c
}

// ID provides a mock function for the type MockRestoreCollTaskView
func (_mock *MockRestoreCollTaskView) ID() string {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for ID")
	}

	var r0 string
	if returnFunc, ok := ret.Get(0).(func() string); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(string)
	}
	return r0
}

// MockRestoreCollTaskView_ID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ID'
type MockRestoreCollTaskView_ID_Call struct {
	*mock.Call
}

// ID is a helper method to define mock.On call
func (_e *MockRestoreCollTaskView_Expecter) ID() *MockRestoreCollTaskView_ID_Call {
	return &MockRestoreCollTaskView_ID_Call{Call: _e.mock.On("ID")}
}

func (_c *MockRestoreCollTaskView_ID_Call) Run(run func()) *MockRestoreCollTaskView_ID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRestoreCollTaskView_ID_Call) Return(s string) *MockRestoreCollTaskView_ID_Call {
	_c.Call.Return(s)
	return _c
}

func (_c *MockRestoreCollTaskView_ID_Call) RunAndReturn(run func() string) *MockRestoreCollTaskView_ID_Call {
	_c.Call.Return(run)
	return _c
}

// Progress provides a mock function for the type MockRestoreCollTaskView
func (_mock *MockRestoreCollTaskView) Progress() int32 {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for Progress")
	}

	var r0 int32
	if returnFunc, ok := ret.Get(0).(func() int32); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(int32)
	}
	return r0
}

// MockRestoreCollTaskView_Progress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Progress'
type MockRestoreCollTaskView_Progress_Call struct {
	*mock.Call
}

// Progress is a helper method to define mock.On call
func (_e *MockRestoreCollTaskView_Expecter) Progress() *MockRestoreCollTaskView_Progress_Call {
	return &MockRestoreCollTaskView_Progress_Call{Call: _e.mock.On("Progress")}
}

func (_c *MockRestoreCollTaskView_Progress_Call) Run(run func()) *MockRestoreCollTaskView_Progress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRestoreCollTaskView_Progress_Call) Return(n int32) *MockRestoreCollTaskView_Progress_Call {
	_c.Call.Return(n)
	return _c
}

func (_c *MockRestoreCollTaskView_Progress_Call) RunAndReturn(run func() int32) *MockRestoreCollTaskView_Progress_Call {
	_c.Call.Return(run)
	return _c
}

// StartTime provides a mock function for the type MockRestoreCollTaskView
func (_mock *MockRestoreCollTaskView) StartTime() time.Time {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for StartTime")
	}

	var r0 time.Time
	if returnFunc, ok := ret.Get(0).(func() time.Time); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(time.Time)
	}
	return r0
}

// MockRestoreCollTaskView_StartTime_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartTime'
type MockRestoreCollTaskView_StartTime_Call struct {
	*mock.Call
}

// StartTime is a helper method to define mock.On call
func (_e *MockRestoreCollTaskView_Expecter) StartTime() *MockRestoreCollTaskView_StartTime_Call {
	return &MockRestoreCollTaskView_StartTime_Call{Call: _e.mock.On("StartTime")}
}

func (_c *MockRestoreCollTaskView_StartTime_Call) Run(run func()) *MockRestoreCollTaskView_StartTime_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRestoreCollTaskView_StartTime_Call) Return(time1 time.Time) *MockRestoreCollTaskView_StartTime_Call {
	_c.Call.Return(time1)
	return _c
}

func (_c *MockRestoreCollTaskView_StartTime_Call) RunAndReturn(run func() time.Time) *MockRestoreCollTaskView_StartTime_Call {
	_c.Call.Return(run)
	return _c
}

// StateCode provides a mock function for the type MockRestoreCollTaskView
func (_mock *MockRestoreCollTaskView) StateCode() backuppb.RestoreTaskStateCode {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for StateCode")
	}

	var r0 backuppb.RestoreTaskStateCode
	if returnFunc, ok := ret.Get(0).(func() backuppb.RestoreTaskStateCode); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(backuppb.RestoreTaskStateCode)
	}
	return r0
}

// MockRestoreCollTaskView_StateCode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StateCode'
type MockRestoreCollTaskView_StateCode_Call struct {
	*mock.Call
}

// StateCode is a helper method to define mock.On call
func (_e *MockRestoreCollTaskView_Expecter) StateCode() *MockRestoreCollTaskView_StateCode_Call {
	return &MockRestoreCollTaskView_StateCode_Call{Call: _e.mock.On("StateCode")}
}

func (_c *MockRestoreCollTaskView_StateCode_Call) Run(run func()) *MockRestoreCollTaskView_StateCode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRestoreCollTaskView_StateCode_Call) Return(restoreTaskStateCode backuppb.RestoreTaskStateCode) *MockRestoreCollTaskView_StateCode_Call {
	_c.Call.Return(restoreTaskStateCode)
	return _c
}

func (_c *MockRestoreCollTaskView_StateCode_Call) RunAndReturn(run func() backuppb.RestoreTaskStateCode) *MockRestoreCollTaskView_StateCode_Call {
	_c.Call.Return(run)
	return _c
}

// TotalSize provides a mock function for the type MockRestoreCollTaskView
func (_mock *MockRestoreCollTaskView) TotalSize() int64 {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for TotalSize")
	}

	var r0 int64
	if returnFunc, ok := ret.Get(0).(func() int64); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Get(0).(int64)
	}
	return r0
}

// MockRestoreCollTaskView_TotalSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TotalSize'
type MockRestoreCollTaskView_TotalSize_Call struct {
	*mock.Call
}

// TotalSize is a helper method to define mock.On call
func (_e *MockRestoreCollTaskView_Expecter) TotalSize() *MockRestoreCollTaskView_TotalSize_Call {
	return &MockRestoreCollTaskView_TotalSize_Call{Call: _e.mock.On("TotalSize")}
}

func (_c *MockRestoreCollTaskView_TotalSize_Call) Run(run func()) *MockRestoreCollTaskView_TotalSize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockRestoreCollTaskView_TotalSize_Call) Return(n int64) *MockRestoreCollTaskView_TotalSize_Call {
	_c.Call.Return(n)
	return _c
}

func (_c *MockRestoreCollTaskView_TotalSize_Call) RunAndReturn(run func() int64) *MockRestoreCollTaskView_TotalSize_Call {
	_c.Call.Return(run)
	return _c
}
