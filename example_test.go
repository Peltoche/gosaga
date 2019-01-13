package gosaga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Peltoche/gosaga/storage"
)

var foo uint
var bar uint

type request struct {
	Amount uint `json:"amount"`
}

func Example() {
	sagaLog := storage.NewMemory()

	foo = 50
	bar = 50

	saga := NewSagaExecutionCoordinator(sagaLog).
		AppendNewSubRequest("debit", debitAction, debitCompensation).
		AppendNewSubRequest("credit", creditAction, creditCompensation)

	saga.StartSaga(context.Background(), json.RawMessage(`{"amount": 10}`))
	// Output:
	// step: _init / done
	// exec: debit
	// Foo 50 -> 40
	// step: debit / done
	// exec: credit
	// Bar 50 -> 60
	// step: credit / done
	// mark saga as done
	// delete saga
}

func Example_with_abort() {
	sagaLog := storage.NewMemory()

	foo = 50
	bar = 50

	saga := NewSagaExecutionCoordinator(sagaLog).
		// The credit step will fail and the debit step will be automatically
		// rollback.
		AppendNewSubRequest("debit", debitAction, debitCompensation).
		AppendNewSubRequest("credit", actionReturningAFailure, creditCompensation)

	saga.StartSaga(context.Background(), json.RawMessage(`{"amount": 10}`))
	// Output:
	// step: _init / done
	// exec: debit
	// Foo 50 -> 40
	// step: debit / done
	// exec: credit
	// failed "debit"
	// revert step: credit / aborted
	// revert : credit
	// Revert Bar 50 -> 60
	// revert step: credit / done
	// revert : debit
	// Revert Foo 40 -> 50
	// revert step: debit / done
	// mark saga as done
	// delete saga
}

func debitAction(ctx context.Context, sagaCtx json.RawMessage) Result {
	var req request
	err := json.Unmarshal(sagaCtx, &req)
	if err != nil {
		return Failure(err, sagaCtx)
	}

	fmt.Printf("Foo %v -> %v\n", foo, foo-req.Amount)
	foo = foo - req.Amount

	return Success(sagaCtx)
}

// WARNING: For the sake of implicity, this Compensation method is not idempotent
// and this is bad. If anything fail between the Printf and the State commit
// this method will be re-run and the use "Foo" will be credited twice.
func debitCompensation(ctx context.Context, sagaCtx json.RawMessage) Result {
	var req request
	err := json.Unmarshal(sagaCtx, &req)
	if err != nil {
		return Failure(err, sagaCtx)
	}

	fmt.Printf("Revert Foo %v -> %v\n", foo, foo+req.Amount)
	foo = foo + req.Amount

	return Success(nil)
}

func creditAction(ctx context.Context, sagaCtx json.RawMessage) Result {
	var req request
	err := json.Unmarshal(sagaCtx, &req)
	if err != nil {
		return Failure(err, sagaCtx)
	}

	fmt.Printf("Bar %v -> %v\n", bar, bar+req.Amount)
	bar = bar + req.Amount

	return Success(sagaCtx)
}

// WARNING: For the sake of implicity, this Compensation method is not idempotent
// and this is bad. If anything fail between the Printf and the State commit
// this method will be re-run and the use "Bar" will be debited twice.
func creditCompensation(ctx context.Context, sagaCtx json.RawMessage) Result {
	var req request
	err := json.Unmarshal(sagaCtx, &req)
	if err != nil {
		return Failure(err, sagaCtx)
	}

	fmt.Printf("Revert Bar %v -> %v\n", bar, bar+req.Amount)
	bar = bar - req.Amount

	return Success(sagaCtx)
}

func actionReturningAFailure(ctx context.Context, sagaCtx json.RawMessage) Result {
	return Failure(errors.New("some-error"), sagaCtx)
}
