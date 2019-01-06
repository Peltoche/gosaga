package gosaga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Peltoche/gosaga/internal/journal"
)

// Journal is an interface used to save all the SEC actions.
//
// It allow to restore its state in case of failure.
type Journal interface {
	CreateNewSaga(ctx context.Context, cmd json.RawMessage) (string, error)
	MarkSagaAsDone(ctx context.Context, sagaID string) error
	DeleteSaga(ctx context.Context, sagaID string)
	MarkSubRequestAsRunning(ctx context.Context, sagaID string, subRequestID string, cmd json.RawMessage) error
	MarkSubRequestAsDone(ctx context.Context, sagaID string, subRequestID string, result json.RawMessage) error
	GetSagaStatus(sagaID string) string
	GetSagaLastEventLog(sagaID string) (string, string, json.RawMessage)
}

// SEC means Saga Execution Coordinator.
//
// It is used to:
// - Interpret and write the Saga Journal.
// - Apply the Saga Sub-Requests
// - Apply the Saga Compensating Sub-Requests when necessary
type SEC struct {
	subRequestDefs subRequestDefs
	journal        Journal
}

// NewSEC instantiate a new Saga Execution Coordinator (SEC).
func NewSEC(storage journal.Storage) *SEC {
	return &SEC{
		subRequestDefs: []subRequestDef{},
		journal:        journal.New(storage),
	}
}

// AppendNewSubRequest append a new SubRequest to the Saga.
func (t *SEC) AppendNewSubRequest(name string, action Action, compensation Action) *SEC {
	t.subRequestDefs = append(t.subRequestDefs, subRequestDef{
		SubRequestID: name,
		Action:       action,
		Compensation: compensation,
	})

	return t
}

// StartSaga create a new Saga saga with the given cmd and run it.
func (t *SEC) StartSaga(ctx context.Context, cmd json.RawMessage) error {
	sagaID, err := t.journal.CreateNewSaga(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to create a new saga: %s", err)
	}

	return t.runSaga(ctx, sagaID)
}

// RunSaga execute the given Saga synchronously.
func (t *SEC) runSaga(ctx context.Context, sagaID string) error {
	for {
		switch t.journal.GetSagaStatus(sagaID) {
		case "running":
			err := t.execNextSubRequestAction(ctx, sagaID)
			if err != nil {
				panic(err)
			}

		case "done":
			fmt.Println("delete saga")
			t.journal.DeleteSaga(ctx, sagaID)
			return nil

		case "abort":
			panic("abort")

		default:
			return fmt.Errorf("unknown saga state: %q", t.journal.GetSagaStatus(sagaID))
		}

	}
}

func (t *SEC) execNextSubRequestAction(ctx context.Context, sagaID string) error {
	step, state, arg := t.journal.GetSagaLastEventLog(sagaID)
	if state == "running" {
		// The previous subRequest is not finished, abort.
		return errors.New("the previous sub-request action/compensation is not finished")
	}
	fmt.Printf("step: %s / %s\n", step, state)

	// Select the next subRequest.
	subReq, err := t.subRequestDefs.GetSubRequestAfter(step)
	if err != nil {
		return fmt.Errorf("failed to select the next sub-request: %s", err)
	}

	if subReq == nil {
		fmt.Println("mark saga as done")
		err = t.journal.MarkSagaAsDone(ctx, sagaID)
		if err != nil {
			return fmt.Errorf("failed to create the saga: %s", err)
		}

		return nil
	}

	fmt.Printf("exec: %s\n", subReq.SubRequestID)
	err = t.journal.MarkSubRequestAsRunning(ctx, sagaID, subReq.SubRequestID, arg)
	if err != nil {
		return fmt.Errorf("failed to create the saga: %s", err)
	}

	result := subReq.Action(ctx, arg)
	if result.IsSuccess() {
		err = t.journal.MarkSubRequestAsDone(ctx, sagaID, subReq.SubRequestID, arg)
		if err != nil {
			return fmt.Errorf("failed to create the saga: %s", err)
		}
	} else {
		panic("should abort")
		//err = t.Journal.MarkSubRequestAsDone(ctx, sagaID, subState.SubRequestID, subState.Arg)
		//if err != nil {
		//return fmt.Errorf("failed to create the saga: %s", err)
		//}
	}

	return nil
}