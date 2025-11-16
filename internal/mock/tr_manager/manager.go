package trmanager

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
)

type MockTrManager struct{}

func (m *MockTrManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (m *MockTrManager) DoWithSettings(ctx context.Context, settings trm.Settings, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func NewMockTrManager() *MockTrManager {
	return &MockTrManager{}
}
