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
func (t *SubRequestMock) Action(ctx context.Context, cmd json.RawMessage) Result {
	args := t.Called(cmd)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(Result)
}

// Compensation Sub-Request mock.
func (t *SubRequestMock) Compensation(ctx context.Context, cmd json.RawMessage) Result {
	args := t.Called(cmd)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(Result)
}

func Test_SEC_StartSaga_success(t *testing.T) {
	cmd := json.RawMessage(`{"key": "value"}`)

	journal := new(journal.Mock)
	subRequest := new(SubRequestMock)
	scheduler := &SEC{subRequestDefs: []subRequestDef{}, journal: journal}
	scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

	// Create and save the Saga
	journal.On("CreateNewSaga", cmd).Return("some-saga-id", nil).Once()

	// There is 3 loops:
	// 1 - Execute "step1"
	// 2 - Mark the saga as done
	// 3 - Delete the saga

	// 1 - Execute the "step1" SubRequest
	journal.On("GetSagaStatus", "some-saga-id").Return("running").Once()
	journal.On("GetSagaLastEventLog", "some-saga-id").Return("_init", "done", cmd).Once()
	journal.On("MarkSubRequestAsRunning", "some-saga-id", "step1", cmd).Return(nil).Once()
	subRequest.On("Action", cmd).Return(Success(cmd)).Once()
	journal.On("MarkSubRequestAsDone", "some-saga-id", "step1", cmd).Return(nil).Once()

	// 2 - Mark the saga as "done"
	journal.On("GetSagaStatus", "some-saga-id").Return("running").Once()
	journal.On("GetSagaLastEventLog", "some-saga-id").Return("step1", "done", cmd).Once()
	journal.On("MarkSagaAsDone", "some-saga-id").Return(nil).Once()

	// 3 - Delete the saga
	journal.On("GetSagaStatus", "some-saga-id").Return("done").Once()
	journal.On("DeleteSaga", "some-saga-id").Once()

	err := scheduler.StartSaga(context.Background(), cmd)
	assert.NoError(t, err)

	journal.AssertExpectations(t)
	subRequest.AssertExpectations(t)
}

func Test_SEC_StartSaga_with_journal_CreateNewSaga_error_should_fail(t *testing.T) {
	cmd := json.RawMessage(`{"key": "value"}`)

	journal := new(journal.Mock)
	subRequest := new(SubRequestMock)
	scheduler := &SEC{subRequestDefs: []subRequestDef{}, journal: journal}
	scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

	// Create and save the Saga
	journal.On("CreateNewSaga", cmd).Return("", errors.New("some-error")).Once()

	err := scheduler.StartSaga(context.Background(), cmd)
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

//func Test_SEC_runSaga_with_an_unknown_state_should_fail(t *testing.T) {
//cmd := json.RawMessage(`{"key": "value"}`)

//journal := new(journal.Mock)
//subRequest := new(SubRequestMock)
//scheduler := &SEC{subRequestDefs: []SubRequestDef{}, journal: journal}
//scheduler.AppendNewSubRequest("step1", subRequest.Action, subRequest.Compensation)

//// Start a saga with an invalid state
//err := scheduler.runSaga(context.Background(), &model.Saga{
//ID:    "some-id",
//State: "invalid-state",
//Cmd:   cmd,
//})

//assert.EqualError(t, err, "unknown saga state: \"invalid-state\"")

//journal.AssertExpectations(t)
//subRequest.AssertExpectations(t)
//}

//func Test_NewSEC_default_generateID_generate_a_valid_uuid(t *testing.T) {
//journal := new(journal.Mock)
//subRequest := new(SubRequestMock)
//scheduler := NewSEC(storage).AppendNewSubRequest("step1", subRequestMock.Action, subRequestMock.Compensation)

//sampleID := scheduler.generateID()

//assert.NotNil(t, uuid.FromStringOrNil(sampleID))

//journal.AssertExpectations(t)
//subRequest.AssertExpectations(t)
//}

//func Test_SEC_executeSubRequest_should_success(t *testing.T) {
//cmd := json.RawMessage(`{"key": "value"}`)

//journal := new(journal.Mock)
//subRequest := new(SubRequestMock)
//scheduler := NewSEC(storage).AppendNewSubRequest("step1", subRequestMock.Action, subRequestMock.Compensation)

//// Create the "start" journal for the sub-request
//journal.On("UpdateSagaState", &model.Saga{
//ID:              "some-id",
//State:           "running",
//SubRequest:      "new-sub-request",
//SubRequestState: "start",
//Cmd:             cmd,
//}).Return(nil).Once()

//subRequest.On("Action", cmd).Once()

//// Create the "done" journal for the sub-request
//journal.On("UpdateSagaState", &model.Saga{
//ID:              "some-id",
//State:           "running",
//SubRequest:      "new-sub-request",
//SubRequestState: "done",
//Cmd:             cmd,
//}).Return(nil).Once()

//err := scheduler.executeSubRequest(context.Background(), &model.Saga{
//ID:              "some-id",
//State:           "running",
//SubRequest:      "some-preview-sub-request",
//SubRequestState: "done",
//Cmd:             cmd,
//}, &SubRequestDef{
//SubRequestID: "new-sub-request",
//Action:       subRequest.Action,
//Compensation: subRequest.Compensation,
//})

//assert.NoError(t, err)

//journal.AssertExpectations(t)
//subRequest.AssertExpectations(t)
//}

//func Test_SEC_executeSubRequest_with_start_write_error_should_fail(t *testing.T) {
//cmd := json.RawMessage(`{"key": "value"}`)

//journal := new(journal.Mock)
//subRequest := new(SubRequestMock)
//scheduler := NewSEC(storage).AppendNewSubRequest("step1", subRequestMock.Action, subRequestMock.Compensation)

//// Create the "start" journal for the sub-request
//journal.On("UpdateSagaState", &model.Saga{
//ID:              "some-id",
//State:           "running",
//SubRequest:      "new-sub-request",
//SubRequestState: "start",
//Cmd:             cmd,
//}).Return(errors.New("some-error")).Once()

//err := scheduler.executeSubRequest(context.Background(), &model.Saga{
//ID:              "some-id",
//State:           "running",
//SubRequest:      "some-preview-sub-request",
//SubRequestState: "done",
//Cmd:             cmd,
//}, &SubRequestDef{
//SubRequestID: "new-sub-request",
//Action:       subRequest.Action,
//Compensation: subRequest.Compensation,
//})

//assert.EqualError(t, err, "some-error")

//journal.AssertExpectations(t)
//subRequest.AssertExpectations(t)
//}

//func Test_SEC_executeSubRequest_with_done_write_error_should_fail(t *testing.T) {
//cmd := json.RawMessage(`{"key": "value"}`)

//journal := new(journal.Mock)
//subRequest := new(SubRequestMock)
//scheduler := NewSEC(storage).AppendNewSubRequest("step1", subRequestMock.Action, subRequestMock.Compensation)

//// Create the "start" journal for the sub-request
//journal.On("UpdateSagaState", &model.Saga{
//ID:              "some-id",
//State:           "running",
//SubRequest:      "new-sub-request",
//SubRequestState: "start",
//Cmd:             cmd,
//}).Return(nil).Once()

//subRequest.On("Action", cmd).Once()

//// Create the "done" journal for the sub-request
//journal.On("UpdateSagaState", &model.Saga{
//ID:              "some-id",
//State:           "running",
//SubRequest:      "new-sub-request",
//SubRequestState: "done",
//Cmd:             cmd,
//}).Return(errors.New("some-error")).Once()

//err := scheduler.executeSubRequest(context.Background(), &model.Saga{
//ID:              "some-id",
//State:           "running",
//SubRequest:      "some-preview-sub-request",
//SubRequestState: "done",
//Cmd:             cmd,
//}, &SubRequestDef{
//SubRequestID: "new-sub-request",
//Action:       subRequest.Action,
//Compensation: subRequest.Compensation,
//})

//assert.EqualError(t, err, "some-error")

//journal.AssertExpectations(t)
//subRequest.AssertExpectations(t)
//}
