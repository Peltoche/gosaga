package storage

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Peltoche/gosaga/model"
)

func Test_Mock_SaveEventLog(t *testing.T) {
	eventlog := new(Mock)

	event := &model.EventLog{
		SagaID: "some-id",
		State:  "some-state",
		Arg:    json.RawMessage(`{"key": "value"}`),
	}

	eventlog.On("SaveEventLog", event).Return(nil)

	eventlog.SaveEventLog(context.Background(), event)

	eventlog.AssertExpectations(t)
}
