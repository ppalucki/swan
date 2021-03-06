package executor

import mock "github.com/stretchr/testify/mock"

// MockExecutor is an autogenerated mock type for the Executor type
type MockExecutor struct {
	mock.Mock
}

// Execute provides a mock function with given fields: command
func (_m *MockExecutor) Execute(command string) (TaskHandle, error) {
	ret := _m.Called(command)

	var r0 TaskHandle
	if rf, ok := ret.Get(0).(func(string) TaskHandle); ok {
		r0 = rf(command)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(TaskHandle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(command)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// String provides a mock function with given fields:
func (_m *MockExecutor) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
