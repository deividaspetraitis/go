package es

import (
	"sync"

	"github.com/deividaspetraitis/go/errors"

	"golang.org/x/exp/slices"
)

// AggregateRoot represents a top entity of the aggregate.
type AggregateRoot struct {
	version Version // version of the aggregate
	state   Events  // state list of events representing current aggregate state
	pending Events  // events are new aggregate events to be handled

	mu sync.Mutex // guard fields above
}

// Root implements es.Aggregate.
func (agg *AggregateRoot) Root() *AggregateRoot {
	return agg
}

// Events returns events that were applied to the aggregate recently.
func (ag *AggregateRoot) Events() []*Event {
	ag.mu.Lock()
	defer ag.mu.Unlock()
	return ag.pending
}

var (
	ErrDuplicateEvent  = errors.New("not valid event: duplicated event")
	ErrVersionMismatch = errors.New("not valid event: version mismatch")
)

// Reply sets given slice of events as persisted state of the aggregate.
// Replied events does not appear in pending state.
// It is not allowed to reply once a call to Apply is made.
func (ar *AggregateRoot) Reply(event []*Event) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	if len(ar.pending) > 0 {
		return errors.New("aggregate is already initialised")
	}

	if !Events(event).Unique() {
		return ErrDuplicateEvent
	}

	if !Events(append(ar.state, event...)).Unique() {
		return ErrDuplicateEvent
	}

	// event version should be larger than last replied event
	last := ar.getLastEvent() // it can be nil
	for _, v := range event {
		if last != nil && last.Version >= v.Version {
			return ErrVersionMismatch
		}

	}

	ar.state = append(ar.state, event...)

	return nil
}

// No events to sync error.
var ErrSyncNoEvents = errors.New("unable to sync an event: no pending events")

// Sync implements Aggregate.
func (ar *AggregateRoot) Sync(event *Event) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	if len(ar.pending) == 0 {
		return ErrSyncNoEvents
	}

	// sync is performed in the sequential order
	// starting with oldest event in the pending state.
	if ar.pending[0].Version != event.Version {
		return ErrVersionMismatch
	}

	ar.state, ar.pending = append(ar.state, ar.pending[0]), ar.pending[1:]

	return nil
}

// getLastEvent returns last aggregates event
// Such event might reside either in pending or state slices.
func (ar *AggregateRoot) getLastEvent() *Event {
	pending, state := len(ar.pending), len(ar.state)
	if state == 0 && pending == 0 {
		return nil
	}

	if pending > 0 {
		return ar.pending[pending-1]
	}

	return ar.state[state-1]
}

// Apply adds event to pending states list.
func (ar *AggregateRoot) Apply(event *Event) error {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	version := append(ar.state.versions(), ar.pending.versions()...)
	if slices.Contains(version, event.Version) {
		return ErrDuplicateEvent
	}

	// current event must be one version greater than last event
	last := ar.getLastEvent() // it can be nil
	if last != nil && last.Version+1 != event.Version {
		return ErrVersionMismatch
	}

	// add to pending list
	ar.pending = append(ar.pending, event)

	return nil
}

// Version returns current Aggregate version.
func (ar *AggregateRoot) Version() Version {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	return ar.version
}

// AdvanceVersion increases Aggregate's version in sequential increasing order.
func (ar *AggregateRoot) AdvanceVersion() Version {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.version += 1
	return ar.version
}

// Aggregate is a single entity capable to build its states based on series of events.
type Aggregate interface {
	// Root returns Aggregate Root.
	Root() *AggregateRoot

	// Sync marks event as persisted and such event will be no longer returned by Events method.
	Sync(event *Event) error

	// Events returns list of recently applied events.
	// Client must call Sync method after persisting new states by this method.
	Events() []*Event

	// Apply applies a new event to the Aggregate.
	// Applied event changes Aggregate state immediately and becomes available under Events method.
	Apply(event *Event) error

	// Reply reconstructs Aggregate state from series of events.
	// Replied events are not returned by Events method.
	Reply(event []*Event) error
}

// ParseAggregateName returns aggregate type name in the string representation.
func ParseAggregateName(v any) string {
	return parseTypeName(v)
}
