# GoSaga

Manage easily your distributed transactions

##  Example

A pretty basic example use to show how it can be used. There is also [a
working example](./example_test.go)

```go
// Debit Action and Compensation
func debitAction(ctx context.Context, cmd json.RawMessage) *gosaga.Result       { /* do a debit */ }
func debitCompensation(ctx context.Context, cmd json.RawMessage) *gosaga.Result { /* revert a debit */ }

// Credit Action and Compensation
func creditAction(ctx context.Context, cmd json.RawMessage) *gosaga.Result       { /* do a credit */ }
func creditCompensation(ctx context.Context, cmd json.RawMessage) *gosaga.Result { /* revert a credit */ }

func main() {
	// Use a in memory storage for testing purpose
	sagaLog := storage.NewMemory()

	saga := gosaga.NewSEC(sagaLog).
		AppendNewSubRequest("debit", debitAction, debitCompensation).
		AppendNewSubRequest("credit", creditAction, creditCompensation)

	saga.StartSaga(context.Background(), json.RawMessage(`{
		"debiter": "foo",
		"crediter": "bar",
		"amount": 10
	}`))
}

```
