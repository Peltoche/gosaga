package gosaga

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetFirstSubRequest_success(t *testing.T) {
	subRequestMock := new(SubRequestMock)

	subRequests := subRequestDefs{
		{
			SubRequestID: "step1",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
		{
			SubRequestID: "step2",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
	}

	res := subRequests.GetFirstSubRequest()

	assert.Equal(t, "step1", res.SubRequestID)
	assert.NotNil(t, res.Action)
	assert.NotNil(t, res.Compensation)

	subRequestMock.AssertExpectations(t)
}

func Test_GetSubRequestDef_success(t *testing.T) {
	subRequestMock := new(SubRequestMock)

	subRequests := subRequestDefs{
		{
			SubRequestID: "step1",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
		{
			SubRequestID: "step2",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
	}

	res := subRequests.GetSubRequestDef("step1")

	assert.Equal(t, "step1", res.SubRequestID)
	assert.NotNil(t, res.Action)
	assert.NotNil(t, res.Compensation)

	subRequestMock.AssertExpectations(t)
}

func Test_GetSubRequestDef_with_an_unknown_id(t *testing.T) {
	subRequestMock := new(SubRequestMock)

	subRequests := subRequestDefs{
		{
			SubRequestID: "step1",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
		{
			SubRequestID: "step2",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
	}

	res := subRequests.GetSubRequestDef("invalid-subrequest-id")

	assert.Nil(t, res)

	subRequestMock.AssertExpectations(t)
}

func Test_GetSubRequestAfter_success(t *testing.T) {
	subRequestMock := new(SubRequestMock)

	subRequests := subRequestDefs{
		{
			SubRequestID: "step1",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
		{
			SubRequestID: "step2",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
	}

	res, err := subRequests.GetSubRequestAfter("step1")

	assert.NoError(t, err)
	assert.Equal(t, "step2", res.SubRequestID)
	assert.NotNil(t, res.Action)
	assert.NotNil(t, res.Compensation)

	subRequestMock.AssertExpectations(t)
}

func Test_GetSubRequestAfter_success_with_no_subrequest_after(t *testing.T) {
	subRequestMock := new(SubRequestMock)

	subRequests := subRequestDefs{
		{
			SubRequestID: "step1",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
		{
			SubRequestID: "step2",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
	}

	res, err := subRequests.GetSubRequestAfter("step2")

	assert.NoError(t, err)
	assert.Nil(t, res)

	subRequestMock.AssertExpectations(t)
}

func Test_GetSubRequestBefore_success(t *testing.T) {
	subRequestMock := new(SubRequestMock)

	subRequests := subRequestDefs{
		{
			SubRequestID: "step1",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
		{
			SubRequestID: "step2",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
	}

	res, err := subRequests.GetSubRequestBefore("step2")

	assert.NoError(t, err)
	assert.Equal(t, "step1", res.SubRequestID)
	assert.NotNil(t, res.Action)
	assert.NotNil(t, res.Compensation)

	subRequestMock.AssertExpectations(t)
}

func Test_GetSubRequestBefore_success_with_no_subrequest_before(t *testing.T) {
	subRequestMock := new(SubRequestMock)

	subRequests := subRequestDefs{
		{
			SubRequestID: "step1",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
		{
			SubRequestID: "step2",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
	}

	res, err := subRequests.GetSubRequestBefore("step1")

	assert.NoError(t, err)
	assert.Nil(t, res)

	subRequestMock.AssertExpectations(t)
}

func Test_GetSubRequestBefore_success_with_an_invalid_subrequest_id(t *testing.T) {
	subRequestMock := new(SubRequestMock)

	subRequests := subRequestDefs{
		{
			SubRequestID: "step1",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
		{
			SubRequestID: "step2",
			Action:       subRequestMock.Action,
			Compensation: subRequestMock.Compensation,
		},
	}

	res, err := subRequests.GetSubRequestBefore("invalid-id")

	assert.EqualError(t, err, `unknown sub-request id "invalid-id"`)
	assert.Nil(t, res)

	subRequestMock.AssertExpectations(t)
}
