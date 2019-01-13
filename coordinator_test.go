package gosaga

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/Peltoche/gosaga/internal/journal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// SubRequestMock is a mock implementation of a Sub-Request and a Compensation
// Request.
type SubRequestMock struct {
	mock.Mock
}

// Action Sub-Request mock.
func (t *SubRequestMock) Action(ctx context.Context, sagaCtx json.RawMessage) Result {
	args := t.Called(sagaCtx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(Result)
}

// Compensation Sub-Request mock.
func (t *SubRequestMock) Compensation(ctx context.Context, sagaCtx json.RawMessage) Result {
	args := t.Called(sagaCtx)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(Result)
}

func Test_SEC_StartSaga_success(t *testing.T) {
	sagaCtx := json.RawMessage(`{"key": "value"}`)

	journal := new(journal.Mock)
	subRequest := new(SubRequestMock)
	scheduler := &SEC{subRequestDefs: []subRequestDef{}, journal: journal}
	scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

	// Create and save the Saga
	journal.On("CreateNewSaga", sagaCtx).Return("some-saga-id", nil).Once()

	// There is 3 loops:
	// 1 - Execute "step1"
	// 2 - Mark the saga as done
	// 3 - Delete the saga

	// 1 - Execute the "step1" SubRequest
	journal.On("GetSagaStatus", "some-saga-id").Return("running").Once()
	journal.On("GetSagaLastEventLog", "some-saga-id").Return("_init", "done", sagaCtx).Once()
	journal.On("MarkSubRequestAsRunning", "some-saga-id", "step1", sagaCtx).Return(nil).Once()
	subRequest.On("Action", sagaCtx).Return(Success(sagaCtx)).Once()
	journal.On("MarkSubRequestAsDone", "some-saga-id", "step1", sagaCtx).Return(nil).Once()

	// 2 - Mark the saga as "done"
	journal.On("GetSagaStatus", "some-saga-id").Return("running").Once()
	journal.On("GetSagaLastEventLog", "some-saga-id").Return("step1", "done", sagaCtx).Once()
	journal.On("MarkSagaAsDone", "some-saga-id").Return(nil).Once()

	// 3 - Delete the saga
	journal.On("GetSagaStatus", "some-saga-id").Return("done").Once()
	journal.On("DeleteSaga", "some-saga-id").Once()

	err := scheduler.StartSaga(context.Background(), sagaCtx)
	assert.NoError(t, err)

	journal.AssertExpectations(t)
	subRequest.AssertExpectations(t)
}

func Test_SEC_StartSaga_with_journal_CreateNewSaga_error_should_fail(t *testing.T) {
	sagaCtx := json.RawMessage(`{"key": "value"}`)

	journal := new(journal.Mock)
	subRequest := new(SubRequestMock)
	scheduler := &SEC{subRequestDefs: []subRequestDef{}, journal: journal}
	scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

	// Create and save the Saga
	journal.On("CreateNewSaga", sagaCtx).Return("", errors.New("some-error")).Once()

	err := scheduler.StartSaga(context.Background(), sagaCtx)
	assert.EqualError(t, err, "failed to create a new saga: some-error")

	journal.AssertExpectations(t)
	subRequest.AssertExpectations(t)
}

func Test_SEC_RunSaga_with_an_invalid_status_should_fail(t *testing.T) {
	journal := new(journal.Mock)
	subRequest := new(SubRequestMock)
	scheduler := &SEC{subRequestDefs: []subRequestDef{}, journal: journal}
	scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

	// GetSagaStatus twice, one time for the switch, one time for the error message.
	journal.On("GetSagaStatus", "some-saga-id").Return("some-invalid-status").Twice()

	err := scheduler.runSaga(context.Background(), "some-saga-id")
	assert.EqualError(t, err, `unknown saga state: "some-invalid-status"`)

	journal.AssertExpectations(t)
	subRequest.AssertExpectations(t)
}

func Test_execNextSubRequestAction_with_a_saga_already_running(t *testing.T) {
	sagaCtx := json.RawMessage(`{"key": "value"}`)

	journal := new(journal.Mock)
	subRequest := new(SubRequestMock)
	scheduler := &SEC{subRequestDefs: []subRequestDef{}, journal: journal}
	scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

	// GetSagaStatus twice, one time for the switch, one time for the error message.
	journal.On("GetSagaLastEventLog", "some-saga-id").Return("step1", "running", sagaCtx).Once()

	err := scheduler.execNextSubRequestAction(context.Background(), "some-saga-id")
	assert.EqualError(t, err, "the previous sub-request action/compensation is not finished")

	journal.AssertExpectations(t)
	subRequest.AssertExpectations(t)
}

func Test_execNextSubRequestAction_with_an_invalid_subrequest_id(t *testing.T) {
	sagaCtx := json.RawMessage(`{"key": "value"}`)

	journal := new(journal.Mock)
	subRequest := new(SubRequestMock)
	scheduler := &SEC{subRequestDefs: []subRequestDef{}, journal: journal}
	scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

	// GetSagaStatus twice, one time for the switch, one time for the error message.
	journal.On("GetSagaLastEventLog", "some-saga-id").Return("unknown-subrequest-id", "done", sagaCtx).Once()

	err := scheduler.execNextSubRequestAction(context.Background(), "some-saga-id")
	assert.EqualError(t, err, `failed to select the next sub-request: unknow sub-request id "unknown-subrequest-id"`)

	journal.AssertExpectations(t)
	subRequest.AssertExpectations(t)
}

func Test_execNextSubRequestAction_with_a_MarkSagaAsDone_error(t *testing.T) {
	sagaCtx := json.RawMessage(`{"key": "value"}`)

	journal := new(journal.Mock)
	subRequest := new(SubRequestMock)
	scheduler := &SEC{subRequestDefs: []subRequestDef{}, journal: journal}
	scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

	// GetSagaStatus twice, one time for the switch, one time for the error message.
	journal.On("GetSagaLastEventLog", "some-saga-id").Return("step1", "done", sagaCtx).Once()
	journal.On("MarkSagaAsDone", "some-saga-id").Return(errors.New("some-error")).Once()

	err := scheduler.execNextSubRequestAction(context.Background(), "some-saga-id")
	assert.EqualError(t, err, "failed to mark the saga as done: some-error")

	journal.AssertExpectations(t)
	subRequest.AssertExpectations(t)
}

func Test_execNextSubRequestAction_with_a_MarkSubRequestAsRunning_error(t *testing.T) {
	sagaCtx := json.RawMessage(`{"key": "value"}`)

	journal := new(journal.Mock)
	subRequest := new(SubRequestMock)
	scheduler := &SEC{subRequestDefs: []subRequestDef{}, journal: journal}
	scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

	// GetSagaStatus twice, one time for the switch, one time for the error message.
	journal.On("GetSagaLastEventLog", "some-saga-id").Return("_init", "done", sagaCtx).Once()
	journal.On("MarkSubRequestAsRunning", "some-saga-id", "step1", sagaCtx).Return(errors.New("some-error")).Once()

	err := scheduler.execNextSubRequestAction(context.Background(), "some-saga-id")
	assert.EqualError(t, err, "failed to mark the subrequest \"step1\" for saga \"some-saga-id\" as running: some-error")

	journal.AssertExpectations(t)
	subRequest.AssertExpectations(t)
}

func Test_execNextSubRequestAction_with_a_MarkSubRequestAsDone_error(t *testing.T) {
	sagaCtx := json.RawMessage(`{"key": "value"}`)

	journal := new(journal.Mock)
	subRequest := new(SubRequestMock)
	scheduler := &SEC{subRequestDefs: []subRequestDef{}, journal: journal}
	scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

	// GetSagaStatus twice, one time for the switch, one time for the error message.
	journal.On("GetSagaLastEventLog", "some-saga-id").Return("_init", "done", sagaCtx).Once()
	journal.On("MarkSubRequestAsRunning", "some-saga-id", "step1", sagaCtx).Return(nil).Once()
	subRequest.On("Action", sagaCtx).Return(Success(sagaCtx)).Once()
	journal.On("MarkSubRequestAsDone", "some-saga-id", "step1", sagaCtx).Return(errors.New("some-error")).Once()

	err := scheduler.execNextSubRequestAction(context.Background(), "some-saga-id")
	assert.EqualError(t, err, "failed to mark the subrequest \"step1\" for saga \"some-saga-id\" as done: some-error")

	journal.AssertExpectations(t)
	subRequest.AssertExpectations(t)
}

func Test_execNextSubRequestAction_with_a_MarkSubRequestAsAborted_error(t *testing.T) {
	sagaCtx := json.RawMessage(`{"key": "value"}`)

	journal := new(journal.Mock)
	subRequest := new(SubRequestMock)
	scheduler := &SEC{subRequestDefs: []subRequestDef{}, journal: journal}
	scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

	// GetSagaStatus twice, one time for the switch, one time for the error message.
	journal.On("GetSagaLastEventLog", "some-saga-id").Return("_init", "done", sagaCtx).Once()
	journal.On("MarkSubRequestAsRunning", "some-saga-id", "step1", sagaCtx).Return(nil).Once()
	subRequest.On("Action", sagaCtx).Return(Failure(errors.New("some-action-error"), sagaCtx)).Once()
	journal.On("MarkSubRequestAsAborted", "some-saga-id", "step1", sagaCtx).Return(errors.New("some-error")).Once()

	err := scheduler.execNextSubRequestAction(context.Background(), "some-saga-id")
	assert.EqualError(t, err, "failed to mark the subrequest \"step1\" for saga \"some-saga-id\" as aborted: some-error")

	journal.AssertExpectations(t)
	subRequest.AssertExpectations(t)
}
