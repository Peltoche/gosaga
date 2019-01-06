package gosaga

import "encoding/json"

// SuccessResponse response returned after a successful Action.
//
// - When it is used as a result for an Action method, the next SubRequest will
//	 be called with the given result parameter. If one of the following SubRequest
//	 fail, the Compensation action will also be called with this result so the
//	 datas inside the result paramter should allows to execute either the next
//	 SubRequest and the current SubRequest Compensation method.
//
// - When it is used as a result for a Compensation method, the next
//   Componsensation method will be called but the result parameter will not be
//	 used.
type SuccessResponse struct {
	status string
	result json.RawMessage
}

// IsSuccess return if the action have be successful.
func (t *SuccessResponse) IsSuccess() bool { return true }

// FailureResponse response returned after a errored Action.
//
// - When it is used as a result for an Action, it change the saga state to
//	 "aborted". An "aborted" saga will trigger all the required Compensation
//	 Actions in order to revert all the change previously make by the Saga.
//
// - When it is used as a result for a Compensation, it will be only logged
//	 and the Compensation method will retry.
type FailureResponse struct {
	status string
	err    error
}

// IsSuccess return if the action have be successful.
func (t *FailureResponse) IsSuccess() bool { return false }

// Success generate a Success response.
func Success(result json.RawMessage) *SuccessResponse {
	return &SuccessResponse{
		status: "success",
		result: result,
	}
}

// Failure return a failure response with the given error.
func Failure(err error) *FailureResponse {
	return &FailureResponse{
		status: "failure",
		err:    err,
	}
}
