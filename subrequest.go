package gosaga

import (
	"context"
	"encoding/json"
	"fmt"
)

// Result of a Sub-Request.
//
// It contains the status and the command for the next Sub-Request.
type Result struct {
	Status string
	Arg    json.RawMessage
}

// Action used for a SubRequest Action or Compensation.
type Action func(ctx context.Context, cmd json.RawMessage) *Result

// SubRequestDef is the definition for an ACID Sub-Request.
type SubRequestDef struct {
	SubRequestID string

	// Action executed by the Sub-Request.
	Action Action

	// Compensation function used to rollback the action in case of failure.
	//
	// **THE COMPENSATION SUBREQUEST NEED TO BE IDEMPOTENT**.
	Compensation Action
}

// SubRequestDefs is the ordered collection of SubRequest.
type subRequestDefs []SubRequestDef

// GetFirstSubRequest return the first Sub-Request to execute.
func (t subRequestDefs) GetFirstSubRequest() *SubRequestDef {
	return &t[0]
}

// GetSubRequestAfte return the next Sub-Request to execute after the given subRequestID.
//
// If there is no more Sub-Request to execute, return nil
func (t subRequestDefs) GetSubRequestAfter(subRequestID string) (*SubRequestDef, error) {
	if subRequestID == "_init" {
		return &t[0], nil
	}

	nextSubReqIDX := -1
	for idx, subReq := range t {
		if subReq.SubRequestID == subRequestID {
			nextSubReqIDX = idx + 1
		}
	}

	if nextSubReqIDX == -1 {
		return nil, fmt.Errorf("unknow sub-request id %q", subRequestID)
	}

	if nextSubReqIDX >= len(t) {
		return nil, nil
	}

	nextSubReq := t[nextSubReqIDX]

	return &nextSubReq, nil
}
