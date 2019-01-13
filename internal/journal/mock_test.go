package journal

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Mock_CreateNewSaga(t *testing.T) {
	mock := new(Mock)

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	mock.On("CreateNewSaga", sagaCtx).Once().Return("some-saga-id", nil)

	sagaID, err := mock.CreateNewSaga(context.Background(), sagaCtx)

	assert.NoError(t, err)
	assert.Equal(t, "some-saga-id", sagaID)

	mock.AssertExpectations(t)
}

func Test_Mock_MarkSubRequestAsRunning(t *testing.T) {
	mock := new(Mock)

	result := json.RawMessage(`{"key": "value"}`)

	mock.On("MarkSubRequestAsRunning", "some-saga-id", "some-subrequest-id", result).Once().Return(nil)

	err := mock.MarkSubRequestAsRunning(context.Background(), "some-saga-id", "some-subrequest-id", result)

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

func Test_Mock_MarkSubRequestAsDone(t *testing.T) {
	mock := new(Mock)

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	mock.On("MarkSubRequestAsDone", "some-saga-id", "some-subrequest-id", sagaCtx).Once().Return(nil)

	err := mock.MarkSubRequestAsDone(context.Background(), "some-saga-id", "some-subrequest-id", sagaCtx)

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

func Test_Mock_MarkSubRequestAsAborted(t *testing.T) {
	mock := new(Mock)

	sagaCtx := json.RawMessage(`{"reason": "some-error"}`)

	mock.On("MarkSubRequestAsAborted", "some-saga-id", "some-subrequest-id", sagaCtx).Once().Return(nil)

	err := mock.MarkSubRequestAsAborted(context.Background(), "some-saga-id", "some-subrequest-id", sagaCtx)

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

func Test_Mock_MarkSagaAsDone(t *testing.T) {
	mock := new(Mock)

	mock.On("MarkSagaAsDone", "some-saga-id").Once().Return(nil)

	err := mock.MarkSagaAsDone(context.Background(), "some-saga-id")

	assert.NoError(t, err)

	mock.AssertExpectations(t)
}

func Test_Mock_GetSagaStatus(t *testing.T) {
	mock := new(Mock)

	mock.On("GetSagaStatus", "some-saga-id").Once().Return("some-status")

	status := mock.GetSagaStatus("some-saga-id")

	assert.Equal(t, "some-status", status)

	mock.AssertExpectations(t)
}

func Test_Mock_GetSagaLastEventLog(t *testing.T) {
	mock := new(Mock)

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	mock.On("GetSagaLastEventLog", "some-saga-id").Once().Return("some-step", "some-state", sagaCtx)

	step, state, res := mock.GetSagaLastEventLog("some-saga-id")

	assert.Equal(t, "some-step", step)
	assert.Equal(t, "some-state", state)
	assert.EqualValues(t, sagaCtx, res)

	mock.AssertExpectations(t)
}

func Test_Mock_GetSagaLastEventLog_with_nil(t *testing.T) {
	mock := new(Mock)

	mock.On("GetSagaLastEventLog", "some-saga-id").Once().Return("", "", nil)

	step, state, res := mock.GetSagaLastEventLog("some-saga-id")

	assert.Empty(t, step)
	assert.Empty(t, state)
	assert.Nil(t, res)

	mock.AssertExpectations(t)
}

func Test_Mock_DeleteSaga(t *testing.T) {
	mock := new(Mock)

	mock.On("DeleteSaga", "some-saga-id").Once().Return()

	mock.DeleteSaga(context.Background(), "some-saga-id")

	mock.AssertExpectations(t)
}
