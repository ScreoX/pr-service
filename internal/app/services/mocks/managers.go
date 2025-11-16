package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type TxManager struct {
	mock.Mock
}

func (m *TxManager) Do(ctx context.Context, operation func(ctx context.Context) error) error {
	args := m.Called(ctx, operation)

	if operation != nil && args.Error(0) == nil {
		return operation(ctx)
	}

	return args.Error(0)
}
