package storage

import (
	"context"

	"github.com/Peltoche/gosaga/model"
	"github.com/stretchr/testify/mock"
)

// Mock implementation of an eventlog storage.
type Mock struct {
	mock.Mock
}

// SaveEventLog mock implementation.
func (t *Mock) SaveEventLog(ctx context.Context, event *model.EventLog) error {
	return t.Called(event).Error(0)
}
