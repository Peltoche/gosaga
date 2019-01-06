package storage

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Peltoche/gosaga/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Memory_SaveEventLog_should_success(t *testing.T) {
	memory := NewMemory()

	event := &model.EventLog{
		SagaID: "some-id",
		State:  "some-state",
		Arg:    json.RawMessage(`{"key": "value"}`),
	}

	err := memory.SaveEventLog(context.Background(), event)
	require.NoError(t, err)

	assert.Len(t, memory.journal, 1)

	assert.EqualValues(t, &memory.journal[0], event)
}
