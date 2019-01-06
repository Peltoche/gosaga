package storage

import (
	"context"
	"sync"

	"github.com/Peltoche/gosaga/model"
)

// Memory eventlog storage using the RAM as storage.
//
// It should be used only for testing purpose as it doesn't ensure any durability
// for the data.
type Memory struct {
	mutex   *sync.Mutex
	journal []model.EventLog
}

// NewMemory instantiate a new Memory.
func NewMemory() *Memory {
	return &Memory{
		mutex:   new(sync.Mutex),
		journal: []model.EventLog{},
	}
}

// SaveEventLog save a new eventlog about a saga Change.
func (t *Memory) SaveEventLog(ctx context.Context, event *model.EventLog) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.journal = append(t.journal, *event)

	return nil
}
