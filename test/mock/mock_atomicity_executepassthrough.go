package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockAtomicExecutorExecutePassthrough struct {
	mock.Mock
}

func NewMockAtomicExecutorExecutePassthrough(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAtomicExecutorExecutePassthrough {
	mock := &MockAtomicExecutorExecutePassthrough{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

func (_m *MockAtomicExecutorExecutePassthrough) Execute(parentCtx context.Context, executeFunc func(context.Context) error) error {
	return executeFunc(parentCtx)
}
