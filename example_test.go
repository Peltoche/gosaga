package gosaga

import (
	"context"
	"encoding/json"
	"errors"

	"log"

	"github.com/Peltoche/gosaga/storage"
)

var foo uint = 50
var bar uint = 50

type request struct {
	Debiter  string `json:"debiter"`
	Crediter string `json:"crediter"`
	Amount   uint   `json:"amount"`
}

func Example() {
	sagaLog := storage.NewMemory()

	saga := NewSEC(sagaLog).
		AppendNewSubRequest("debit", debitAction, debitCompensation).
		AppendNewSubRequest("credit", creditAction, creditCompensation)

	saga.StartSaga(context.Background(), json.RawMessage(`{
		"debiter": "foo",
		"crediter": "bar",
		"amount": 10
	}`))
	// Output:
	// step: _init / done
	// exec: debit
	// step: debit / done
	// exec: credit
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

	switch req.Debiter {
	case "foo":
		log.Printf("Foo %v -> %v\n", foo, foo-req.Amount)
		foo = foo - req.Amount
	case "bar":
		log.Printf("Bar %v -> %v\n", bar, bar-req.Amount)
		bar = bar - req.Amount
	default:
		return Failure(errors.New("unknown target"))
	}

	return Success(cmd)
}

func debitCompensation(ctx context.Context, cmd json.RawMessage) Result {
	var req request
	err := json.Unmarshal(cmd, &req)
	if err != nil {
		return Failure(err)
	}

	switch req.Debiter {
	case "foo":
		log.Printf("Revert Foo %v -> %v\n", foo, foo+req.Amount)
		foo = foo + req.Amount
	case "bar":
		log.Printf("Revert Bar %v -> %v\n", bar, bar+req.Amount)
		bar = bar + req.Amount
	default:
		return Failure(errors.New("unknown target"))
	}

	return Success(cmd)
}

func creditAction(ctx context.Context, cmd json.RawMessage) Result {
	var req request
	err := json.Unmarshal(cmd, &req)
	if err != nil {
		return Failure(err)
	}

	switch req.Crediter {
	case "foo":
		log.Printf("Foo %v -> %v\n", foo, foo+req.Amount)
		foo = foo + req.Amount
	case "bar":
		log.Printf("Bar %v -> %v\n", bar, bar+req.Amount)
		bar = bar + req.Amount
	default:
		return Failure(errors.New("unknown target"))
	}

	return Success(cmd)
}

func creditCompensation(ctx context.Context, cmd json.RawMessage) Result {
	var req request
	err := json.Unmarshal(cmd, &req)
	if err != nil {
		return Failure(err)
	}

	switch req.Crediter {
	case "foo":
		log.Printf("Revert Foo %v -> %v\n", foo, foo+req.Amount)
		foo = foo - req.Amount
	case "bar":
		log.Printf("Revert Bar %v -> %v\n", bar, bar+req.Amount)
		bar = bar - req.Amount
	default:
		return Failure(errors.New("unknown target"))
	}

	return Success(cmd)
}
