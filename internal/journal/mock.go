package journal

import (
	"context"
	"encoding/json"

	"github.com/stretchr/testify/mock"
)

// Mock implementation of a Journal.
type Mock struct {
	mock.Mock
}

// CreateNewSaga mock.
func (t *Mock) CreateNewSaga(ctx context.Context, cmd json.RawMessage) (string, error) {
	args := t.Called(cmd)

	return args.String(0), args.Error(1)
}

// MarkSubRequestAsRunning mock.
func (t *Mock) MarkSubRequestAsRunning(ctx context.Context, sagaID string, subRequestID string, cmd json.RawMessage) error {
	return t.Called(sagaID, subRequestID, cmd).Error(0)
}

// MarkSubRequestAsDone mock.
func (t *Mock) MarkSubRequestAsDone(ctx context.Context, sagaID string, subRequestID string, result json.RawMessage) error {
	return t.Called(sagaID, subRequestID, result).Error(0)
}

// MarkSagaAsDone mock.
func (t *Mock) MarkSagaAsDone(ctx context.Context, sagaID string) error {
	return t.Called(sagaID).Error(0)
}

// GetSagaStatus mock.
func (t *Mock) GetSagaStatus(sagaID string) string {
	return t.Called(sagaID).String(0)
}

// GetSagaLastEventLog mock.
func (t *Mock) GetSagaLastEventLog(sagaID string) (string, string, json.RawMessage) {
	args := t.Called(sagaID)

	if args.Get(2) == nil {
		return "", "", nil
	}

	return args.String(0), args.String(1), args.Get(2).(json.RawMessage)
}

// DeleteSaga mock.
func (t *Mock) DeleteSaga(ctx context.Context, sagaID string) {
	t.Called(sagaID)
}
