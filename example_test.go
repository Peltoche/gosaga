package gosaga

import (
	"context"
	"encoding/json"
	"fmt"

	"log"

	"github.com/Peltoche/gosaga/storage"
)

var foo uint = 50
var bar uint = 50

type request struct {
	Amount uint `json:"amount"`
}

func Example() {
	sagaLog := storage.NewMemory()

	saga := NewSEC(sagaLog).
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

func debitAction(ctx context.Context, cmd json.RawMessage) Result {
	var req request
	err := json.Unmarshal(cmd, &req)
	if err != nil {
		return Failure(err)
	}

	foo = foo - req.Amount
	fmt.Printf("Foo %v -> %v\n", foo, foo-req.Amount)

	return Success(cmd)
}

// WARNING: For the sake of implicity, this Compensation method is not idempotent
// and this is bad. If anything fail between the Printf and the State commit
// this method will be re-run and the use "Foo" will be credited twice.
func debitCompensation(ctx context.Context, cmd json.RawMessage) Result {
	var req request
	err := json.Unmarshal(cmd, &req)
	if err != nil {
		return Failure(err)
	}

	foo = foo + req.Amount
	fmt.Printf("Revert Foo %v -> %v\n", foo, foo+req.Amount)

	return Success(nil)
}

func creditAction(ctx context.Context, cmd json.RawMessage) Result {
	var req request
	err := json.Unmarshal(cmd, &req)
	if err != nil {
		return Failure(err)
	}

	bar = bar + req.Amount
	fmt.Printf("Bar %v -> %v\n", bar, bar+req.Amount)

	return Success(cmd)
}

// WARNING: For the sake of implicity, this Compensation method is not idempotent
// and this is bad. If anything fail between the Printf and the State commit
// this method will be re-run and the use "Bar" will be debited twice.
func creditCompensation(ctx context.Context, cmd json.RawMessage) Result {
	var req request
	err := json.Unmarshal(cmd, &req)
	if err != nil {
		return Failure(err)
	}

	bar = bar - req.Amount
	log.Printf("Revert Bar %v -> %v\n", bar, bar+req.Amount)

	return Success(cmd)
}
