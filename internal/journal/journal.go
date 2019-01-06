package journal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Peltoche/gosaga/model"
	uuid "github.com/satori/go.uuid"
)

// Storage is the driver used to save the eventlogs in a persistent way.
type Storage interface {
	SaveEventLog(ctx context.Context, state *model.EventLog) error
}

// Journal handle all the interfactions with the eventlogs.
//
// It contains an internal map which contains all the eventslogs by Saga.
type Journal struct {
	storage    Storage
	journal    map[string]model.Saga
	generateID func() string
}

// New instanciate a new Journal.
func New(storage Storage) *Journal {
	return &Journal{
		storage:    storage,
		journal:    map[string]model.Saga{},
		generateID: func() string { return uuid.NewV4().String() },
	}
}

// CreateNewSaga mark the given Saga a started.
func (t *Journal) CreateNewSaga(ctx context.Context, sagaCtx json.RawMessage) (string, error) {
	sagaID := t.generateID()

	eventLog := model.EventLog{SagaID: sagaID, Step: "_init", State: "done", Context: sagaCtx}
	err := t.storage.SaveEventLog(ctx, &eventLog)
	if err != nil {
		return "", fmt.Errorf("failed to save into the storage: %s", err)
	}

	t.journal[sagaID] = model.Saga{
		ID:        sagaID,
		Status:    "running",
		EventLogs: []model.EventLog{eventLog},
	}

	return sagaID, nil
}

// MarkSubRequestAsRunning make the given Sub-Request as started for the given Saga.
func (t *Journal) MarkSubRequestAsRunning(ctx context.Context, sagaID string, subRequestID string, sagaCtx json.RawMessage) error {
	saga, ok := t.journal[sagaID]
	if !ok {
		return fmt.Errorf("saga %q not found into the journal", sagaID)
	}

	eventLog := model.EventLog{SagaID: sagaID, Step: subRequestID, State: "running", Context: sagaCtx}
	err := t.storage.SaveEventLog(ctx, &eventLog)
	if err != nil {
		return fmt.Errorf("failed to save into the storage: %s", err)
	}

	saga.EventLogs = append(saga.EventLogs, eventLog)

	t.journal[sagaID] = saga

	return nil
}

// MarkSubRequestAsDone make the given Sub-Request as started for the given Saga.
func (t *Journal) MarkSubRequestAsDone(ctx context.Context, sagaID string, subRequestID string, sagaCtx json.RawMessage) error {
	saga, ok := t.journal[sagaID]
	if !ok {
		return fmt.Errorf("saga %q not found into the journal", sagaID)
	}

	subRequestCurrentStep := ""
	for _, eventLog := range saga.EventLogs {
		if strings.HasPrefix(eventLog.Step, subRequestID) {
			subRequestCurrentStep = eventLog.State
		}
	}

	if subRequestCurrentStep == "" {
		return errors.New("expected current state to be \"running\", have not previous state")
	}

	if subRequestCurrentStep != "running" {
		return fmt.Errorf("expected current state to be \"running\", have %q", subRequestCurrentStep)
	}

	eventLog := model.EventLog{SagaID: sagaID, Step: subRequestID, State: "done", Context: sagaCtx}
	err := t.storage.SaveEventLog(ctx, &eventLog)
	if err != nil {
		return fmt.Errorf("failed to save into the storage: %s", err)
	}

	saga.EventLogs = append(saga.EventLogs, eventLog)

	t.journal[sagaID] = saga

	return nil
}

// MarkSubRequestAsAborted make the given Sub-Request and saga as aborted for the given Saga.
func (t *Journal) MarkSubRequestAsAborted(ctx context.Context, sagaID string, subRequestID string, sagaCtx json.RawMessage) error {
	saga, ok := t.journal[sagaID]
	if !ok {
		return fmt.Errorf("saga %q not found into the journal", sagaID)
	}

	subRequestCurrentStep := ""
	for _, eventLog := range saga.EventLogs {
		if strings.HasPrefix(eventLog.Step, subRequestID) {
			subRequestCurrentStep = eventLog.State
		}
	}

	if subRequestCurrentStep == "" {
		return errors.New("expected current state to be \"running\", have not previous state")
	}

	if subRequestCurrentStep != "running" {
		return fmt.Errorf("expected current state to be \"running\", have %q", subRequestCurrentStep)
	}

	eventLog := model.EventLog{SagaID: sagaID, Step: subRequestID, State: "aborted", Context: sagaCtx}
	err := t.storage.SaveEventLog(ctx, &eventLog)
	if err != nil {
		return fmt.Errorf("failed to save into the storage: %s", err)
	}

	saga.Status = "aborted"
	saga.EventLogs = append(saga.EventLogs, eventLog)

	t.journal[sagaID] = saga

	return nil
}

// MarkSagaAsDone mark the given Saga a done.
func (t *Journal) MarkSagaAsDone(ctx context.Context, sagaID string) error {
	saga, ok := t.journal[sagaID]
	if !ok {
		return fmt.Errorf("saga %q not found into the journal", sagaID)
	}

	subRequestCurrentStep := saga.EventLogs[len(saga.EventLogs)-1].State
	if subRequestCurrentStep != "done" {
		return fmt.Errorf("expected current state to be \"done\", have %q", subRequestCurrentStep)
	}
	err := t.storage.SaveEventLog(ctx, &model.EventLog{SagaID: sagaID, Step: "_finish", State: "done"})
	if err != nil {
		return fmt.Errorf("failed to save into the storage: %s", err)
	}

	saga.Status = "done"

	t.journal[sagaID] = saga

	return nil
}

// DeleteSaga remove the saga from the local journal but keep it into the storage.
func (t *Journal) DeleteSaga(ctx context.Context, sagaID string) {
	delete(t.journal, sagaID)
}

// GetSagaStatus return the status for the given sagaID.
func (t *Journal) GetSagaStatus(sagaID string) string {
	saga, exists := t.journal[sagaID]

	if !exists {
		return ""
	}

	return saga.Status
}

// GetSagaLastEventLog return the last eventlog for a given saga.
func (t *Journal) GetSagaLastEventLog(sagaID string) (string, string, json.RawMessage) {
	saga, exists := t.journal[sagaID]

	if !exists || len(saga.EventLogs) == 0 {
		return "", "", nil
	}

	// Retrieve the last save State
	eventLog := saga.EventLogs[len(saga.EventLogs)-1]

	return eventLog.Step, eventLog.State, eventLog.Context
}
