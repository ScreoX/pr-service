package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type TimeProvider struct {
	mock.Mock
}

func (m *TimeProvider) Now() time.Time {
	args := m.Called()

	return args.Get(0).(time.Time)
}

type RandomProvider struct {
	mock.Mock
}

func (m *RandomProvider) Shuffle(n int, swap func(i, j int)) {
	m.Called(n, swap)
}

func (m *RandomProvider) Intn(n int) int {
	args := m.Called(n)

	return args.Int(0)
}
