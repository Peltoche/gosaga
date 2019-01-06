package gosaga

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Success(t *testing.T) {
	result := json.RawMessage(`{"key": "value"}`)
	res := Success(result)

	assert.True(t, res.IsSuccess())
	assert.EqualValues(t, &SuccessResponse{
		status: "success",
		result: result,
	}, res)
}

func Test_Failure(t *testing.T) {
	err := errors.New("some-error")
	res := Failure(err)

	assert.False(t, res.IsSuccess())
	assert.EqualValues(t, &FailureResponse{
		status: "failure",
		err:    err,
	}, res)
}
