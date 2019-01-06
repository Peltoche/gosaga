package journal

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Peltoche/gosaga/model"
	"github.com/Peltoche/gosaga/storage"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_New_default_generateID_method(t *testing.T) {
	journal := New(nil)

	assert.NotNil(t, uuid.FromStringOrNil(journal.generateID()))
}

func Test_Journal_CreateNewSaga_success(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)

	id, err := journal.CreateNewSaga(context.Background(), sagaCtx)

	assert.NoError(t, err)
	assert.Equal(t, "some-saga-id", id)

	storageMock.AssertExpectations(t)
}

func Test_Journal_CreateNewSaga_with_a_storage_error(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(errors.New("some-error"))

	id, err := journal.CreateNewSaga(context.Background(), sagaCtx)

	assert.EqualError(t, err, `failed to save into the storage: some-error`)
	assert.Empty(t, id)

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsRunning_success(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as running
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "running", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsRunning(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsRunning_with_storage_error(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as running
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "running", Context: sagaCtx}).Once().Return(errors.New("some-error"))
	err = journal.MarkSubRequestAsRunning(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	assert.EqualError(t, err, "failed to save into the storage: some-error")

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsRunning_with_an_unknown_saga(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	err := journal.MarkSubRequestAsRunning(context.Background(), "some-unknown-saga-id", "some-subrequest-id", sagaCtx)

	assert.EqualError(t, err, "saga \"some-unknown-saga-id\" not found into the journal")

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsDone_success(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as running
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "running", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsRunning(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	// Mark the subrequest as done
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "done", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsDone(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsDone_with_storage_error(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as running
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "running", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsRunning(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	// Mark the subrequest as done
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "done", Context: sagaCtx}).Once().Return(errors.New("some-error"))
	err = journal.MarkSubRequestAsDone(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	assert.EqualError(t, err, "failed to save into the storage: some-error")

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsDone_with_an_unknown_saga(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	err := journal.MarkSubRequestAsDone(context.Background(), "some-unknown-saga-id", "some-subrequest-id", sagaCtx)

	assert.EqualError(t, err, "saga \"some-unknown-saga-id\" not found into the journal")

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsDone_with_the_subrequest_not_in_running_state(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as running
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "running", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsRunning(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	// Mark the subrequest as done
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "done", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsDone(context.Background(), sagaID, "some-subrequest-id", sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as done AGAIN. It should fail as the subrequest is in the "done" State.
	err = journal.MarkSubRequestAsDone(context.Background(), sagaID, "some-subrequest-id", sagaCtx)
	assert.EqualError(t, err, `expected current state to be "running", have "done"`)

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsDone_with_not_subrequest_previous_state(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as done without calling the MarkSubRequestAsRunning before.
	err = journal.MarkSubRequestAsDone(context.Background(), sagaID, "some-subrequest-id", sagaCtx)
	assert.EqualError(t, err, "expected current state to be \"running\", have not previous state")

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSagaAsDone_success(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as done without calling the MarkSubRequestAsRunning before.
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_finish", State: "done"}).Once().Return(nil)
	err = journal.MarkSagaAsDone(context.Background(), sagaID)

	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSagaAsDone_with_unknown_sagaID(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	// Mark an unknown saga as done.
	err := journal.MarkSagaAsDone(context.Background(), "some-invalid-saga-id")

	assert.EqualError(t, err, `saga "some-invalid-saga-id" not found into the journal`)

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSagaAsDone_with_a_running_subrequest(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as running
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "running", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsRunning(context.Background(), sagaID, "some-subrequest-id", sagaCtx)
	require.NoError(t, err)

	// Mark the saga as "done".
	err = journal.MarkSagaAsDone(context.Background(), sagaID)

	assert.EqualError(t, err, `expected current state to be "done", have "running"`)

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSagaAsDone_with_a_storage_error(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the saga as "done".
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_finish", State: "done"}).Once().Return(errors.New("some-error"))
	err = journal.MarkSagaAsDone(context.Background(), sagaID)

	assert.EqualError(t, err, `failed to save into the storage: some-error`)

	storageMock.AssertExpectations(t)
}

func Test_Journal_GetSagaStatus_success(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	status := journal.GetSagaStatus(sagaID)

	assert.Equal(t, "running", status)

	storageMock.AssertExpectations(t)
}

func Test_Journal_GetSagaStatus_with_unknown_sagaID(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	status := journal.GetSagaStatus("some-invalid-saga-id")

	assert.Empty(t, status)

	storageMock.AssertExpectations(t)
}

func Test_Journal_GetSagaLastEventLog_success(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	step, state, arg := journal.GetSagaLastEventLog(sagaID)

	assert.Equal(t, "_init", step)
	assert.Equal(t, "done", state)
	assert.EqualValues(t, sagaCtx, arg)

	storageMock.AssertExpectations(t)
}

func Test_Journal_GetSagaLastEventLog_with_an_unknown_sagaID(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	step, state, arg := journal.GetSagaLastEventLog("some-unknown-saga-id")

	assert.Empty(t, step)
	assert.Empty(t, state)
	assert.Empty(t, arg)

	storageMock.AssertExpectations(t)
}

func Test_Journal_DeleteSaga_success(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	journal.DeleteSaga(context.Background(), sagaID)
	assert.Len(t, journal.journal, 0)

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsAborted_success(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as running
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "running", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsRunning(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	// Mark the subrequest as aborted
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "aborted", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsAborted(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsAborted_with_storage_error(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as running
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "running", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsRunning(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	// Mark the subrequest as aborted
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "aborted", Context: sagaCtx}).Once().Return(errors.New("some-error"))
	err = journal.MarkSubRequestAsAborted(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	assert.EqualError(t, err, "failed to save into the storage: some-error")

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsAborted_with_an_unknown_saga(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	err := journal.MarkSubRequestAsAborted(context.Background(), "some-unknown-saga-id", "some-subrequest-id", sagaCtx)

	assert.EqualError(t, err, "saga \"some-unknown-saga-id\" not found into the journal")

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsAborted_with_the_subrequest_not_in_running_state(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as running
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "running", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsRunning(context.Background(), sagaID, "some-subrequest-id", sagaCtx)

	// Mark the subrequest as done
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "some-subrequest-id", State: "done", Context: sagaCtx}).Once().Return(nil)
	err = journal.MarkSubRequestAsDone(context.Background(), sagaID, "some-subrequest-id", sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as aborted. It should fail as the subrequest is in not in the "running" State.
	err = journal.MarkSubRequestAsAborted(context.Background(), sagaID, "some-subrequest-id", sagaCtx)
	assert.EqualError(t, err, `expected current state to be "running", have "done"`)

	storageMock.AssertExpectations(t)
}

func Test_Journal_MarkSubRequestAsAborted_with_not_subrequest_previous_state(t *testing.T) {
	storageMock := new(storage.Mock)
	journal := New(storageMock)
	journal.generateID = func() string { return "some-saga-id" }

	sagaCtx := json.RawMessage(`{"key": "value"}`)

	// Initialize the saga
	storageMock.On("SaveEventLog", &model.EventLog{SagaID: "some-saga-id", Step: "_init", State: "done", Context: sagaCtx}).Once().Return(nil)
	sagaID, err := journal.CreateNewSaga(context.Background(), sagaCtx)
	require.NoError(t, err)

	// Mark the subrequest as done without calling the MarkSubRequestAsRunning before.
	err = journal.MarkSubRequestAsAborted(context.Background(), sagaID, "some-subrequest-id", sagaCtx)
	assert.EqualError(t, err, "expected current state to be \"running\", have not previous state")

	storageMock.AssertExpectations(t)
}
