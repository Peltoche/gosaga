package model

import (
	"encoding/json"
)

// Saga represent a distributed transaction.
type Saga struct {
	ID        string
	Status    string
	EventLogs []EventLog
}

// EventLog log a change into a Saga state.
type EventLog struct {
	SagaID string
	Step   string
	State  string
	Arg    json.RawMessage
}
