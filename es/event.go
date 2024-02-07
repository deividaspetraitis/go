package es

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/deividaspetraitis/go/errors"
	"github.com/deividaspetraitis/go/slices"
)

// registered events for the aggregates
var (
	events     map[string]map[string]EventRegisterFunc
	aggregates map[string]Aggregate
	mu         sync.Mutex
)

// EventRegisterFunc is used to register an event along aggregate.
type EventRegisterFunc func() MarshalUnmarshaler

// init initialises package state by initialising variables
// for registering supported aggregates and their events.
func init() {
	events = make(map[string]map[string]EventRegisterFunc)
	aggregates = make(map[string]Aggregate)
}

// Version is the event version.
type Version uint64

// Event represents single aggregate state.
// Sequence of events represents sequence of different states in time
// used to restore aggregate to its most recent state.
type Event struct {
	AggregateID string             // Aggregate ID
	Version     Version            // Event version
	Type        string             // Name of Data type
	Aggregate   any                // Aggregrate type
	Timestamp   time.Time          // Event creation time
	Data        MarshalUnmarshaler // Actual event data type
	Metadata    []byte             // Additional data type
}

// NewEvent constructs and returns new Event based on provided inputs.
func NewEvent(id string, agg Aggregate, event MarshalUnmarshaler) *Event {
	return &Event{
		AggregateID: id,
		Version:     agg.Root().AdvanceVersion(), // increase version
		Aggregate:   agg,
		Type:        ParseEventName(event),
		Timestamp:   time.Now().UTC(),
		Data:        event,
	}
}

type MarshalUnmarshaler interface {
	json.Marshaler
	json.Unmarshaler
}

// Events represents a slice of pointers to Event.
type Events []*Event

// Unique reports whether events are unique based on their version or not.
func (e Events) Unique() bool {
	versions := e.versions()
	return len(slices.Unique(versions)) == len(versions)
}

// parseVersions returns versions of all events.
func (e Events) versions() []Version {
	var version []Version
	for _, v := range e {
		version = append(version, v.Version)
	}
	return version
}

// RegisterAggregateEvent registers Aggregate supported event.
func RegisterAggregateEvent(agg Aggregate, reg EventRegisterFunc) {
	mu.Lock()
	aggname := ParseAggregateName(agg)
	eventname := ParseEventName(reg())
	aggregates[aggname] = agg
	if events[aggname] == nil {
		events[aggname] = make(map[string]EventRegisterFunc)
	}
	events[aggname][eventname] = reg
	mu.Unlock()
}

// GetAggregateEvent returns registered Aggregate event by its name.
func GetAggregateEvent(aggregate Aggregate, event string) (MarshalUnmarshaler, error) {
	mu.Lock()
	defer mu.Unlock()

	aggname := ParseAggregateName(aggregate)
	v, ok := events[aggname][event]
	if !ok {
		return nil, errors.New("not found")
	}

	return v(), nil
}

// ParseEventName returns event type name in the string representation.
func ParseEventName(v any) string {
	return parseTypeName(v)
}
