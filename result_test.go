package gosaga

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Success(t *testing.T) {
	context := json.RawMessage(`{"key": "value"}`)
	res := Success(context)

	assert.True(t, res.IsSuccess())
	assert.EqualValues(t, &SuccessResponse{
		status:  "success",
		context: context,
	}, res)
}

func Test_Failure(t *testing.T) {
	context := json.RawMessage(`{"key": "value"}`)
	err := errors.New("some-error")
	res := Failure(err, context)

	assert.False(t, res.IsSuccess())
	assert.EqualValues(t, &FailureResponse{
		status:  "failure",
		context: context,
		err:     err,
	}, res)
}
