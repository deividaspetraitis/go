package es

import (
	"testing"

	"github.com/deividaspetraitis/go/errors"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testEntityInitialized struct {
	ID   string
	Name string
}

// Implements es.MarshalUnmarshaler
func (t *testEntityInitialized) UnmarshalJSON(b []byte) error {
	return nil
}

// Implements es.MarshalUnmarshaler
func (t *testEntityInitialized) MarshalJSON() ([]byte, error) {
	return nil, nil
}

type testIncCounter struct {
	By int
}

// Implements es.MarshalUnmarshaler
func (t *testIncCounter) UnmarshalJSON(b []byte) error {
	return nil
}

// Implements es.MarshalUnmarshaler
func (t *testIncCounter) MarshalJSON() ([]byte, error) {
	return nil, nil
}

type testDecCounter struct {
	By int
}

// Implements es.MarshalUnmarshaler
func (t *testDecCounter) UnmarshalJSON(b []byte) error {
	return nil
}

// Implements es.MarshalUnmarshaler
func (t *testDecCounter) MarshalJSON() ([]byte, error) {
	return nil, nil
}

type testEntity struct {
	ID      string
	Name    string
	Counter int
}

func (e *testEntity) on(event *Event) error {
	switch event := event.Data.(type) {
	case *testEntityInitialized:
		e.ID = event.ID
		e.Name = event.Name
	case *testDecCounter:
		e.Counter -= event.By
	case *testIncCounter:
		e.Counter += event.By
	default:
		return errors.Newf("unsupported event: %#v", e)
	}

	return nil
}

type testAggregate struct {
	AggregateRoot
	testEntity
}

func (a *testAggregate) Apply(event *Event) error {
	if err := a.AggregateRoot.Apply(event); err != nil {
		return err
	}
	return a.on(event)
}

// TestApply verifies Aggregate Reply method behavior.
func TestApply(t *testing.T) {
	var aggregate testAggregate

	first := NewEvent("1", &aggregate, &testEntityInitialized{
		ID:   "1",
		Name: "test",
	})

	NewEvent("1", &aggregate, &testIncCounter{
		By: 10,
	})

	if err := (&aggregate).Apply(first); err != nil {
		t.Fatalf("#%d got %v, want %v", 0, err, nil)
	}

	// Apply should adjust pending
	if !cmp.Equal(Events{first}, aggregate.pending, cmpopts.IgnoreUnexported(testAggregate{}, AggregateRoot{})) {
		t.Errorf("#%d got %v, want %v", 0, Events{first}, aggregate.pending)
	}

	// Apply should not effect state
	if want := 0; len(aggregate.state) != want {
		t.Errorf("#%d got %v, want %v", 1, len(aggregate.state), want)
	}

	// Apply should not allow to apply already applied events
	if err := (&aggregate).Apply(first); !errors.Is(err, ErrDuplicateEvent) {
		t.Fatalf("got %v, want %v", err, ErrDuplicateEvent)
	}

	// Apply should not allow to reply event with lower version
	if err := (&aggregate).Apply(&Event{
		AggregateID: "1",
		Version:     0, // fake version
		Aggregate:   &aggregate,
	}); !errors.Is(err, ErrVersionMismatch) {
		t.Fatalf("got %v, want %v", err, ErrVersionMismatch)
	}
}

// TestReply verifies Aggregate Reply method behavior.
func TestReply(t *testing.T) {
	var aggregate testAggregate

	first := NewEvent("1", &aggregate, &testEntityInitialized{
		ID:   "1",
		Name: "test",
	})

	if err := (&aggregate).Reply([]*Event{
		first,
	}); err != nil {
		t.Fatalf("#%d got %v, want %v", 0, err, nil)
	}

	// Reply should adjust state
	if !cmp.Equal(Events{first}, aggregate.state, cmpopts.IgnoreUnexported(testAggregate{}, AggregateRoot{})) {
		t.Errorf("#%d got %v, want %v", 0, Events{first}, aggregate.state)
	}

	// Reply should not effect pending
	if want := 0; len(aggregate.pending) != want {
		t.Errorf("#%d got %v, want %v", 1, len(aggregate.pending), want)
	}

	// Reply should not allow to reply already replied events
	if err := (&aggregate).Reply([]*Event{
		first,
	}); !errors.Is(err, ErrDuplicateEvent) {
		t.Fatalf("got %v, want %v", err, ErrDuplicateEvent)
	}

	// Reply should not allow to reply event with lower version
	if err := (&aggregate).Reply([]*Event{
		{
			AggregateID: "1",
			Version:     0, // fake version
			Aggregate:   &aggregate,
		},
	}); !errors.Is(err, ErrVersionMismatch) {
		t.Fatalf("got %v, want %v", err, ErrVersionMismatch)
	}
}

// TestSync verifies Aggregate Sync method behavior.
func TestSync(t *testing.T) {
	var aggregate testAggregate

	events := []*Event{
		NewEvent("1", &aggregate, &testIncCounter{
			By: 10,
		}),
		NewEvent("1", &aggregate, &testDecCounter{
			By: 10,
		}),
		NewEvent("1", &aggregate, &testIncCounter{
			By: 10,
		}),
	}

	// iterate over the events and apply.
	// aggregate.pending should reflect this operation.
	for i, v := range events {
		// Apply should not result in error
		if err := (&aggregate).Apply(v); err != nil {
			t.Fatalf("#%d got %v, want %v", i, err, nil)
		}

		// expected slice grows as we go through events.
		expected := Events(events[0 : i+1])
		if !cmp.Equal(expected, aggregate.pending, cmpopts.IgnoreUnexported(testAggregate{}, AggregateRoot{})) {
			t.Errorf("#%d got %v, want %v", i, aggregate.pending, expected)
		}

		// Apply should not effect state
		if want := 0; len(aggregate.state) != want {
			t.Errorf("#%d got %v, want %v", 1, len(aggregate.state), want)
		}
	}

	// sync event that does not exist or otherwise is in wrong order
	// should result in version mismatch error
	for i := len(events) - 1; i >= 1; i-- {
		if err := (&aggregate).Sync(events[i]); !errors.Is(err, ErrVersionMismatch) {
			t.Fatalf("#%d got nil error, want %v", i, ErrVersionMismatch)
		}
	}

	// iterate over the events and sync each one.
	// aggregate.pending should decrease
	// aggregate.state should increase
	for i, v := range events {
		// Sync should not result in error
		if err := (&aggregate).Sync(v); err != nil {
			t.Fatalf("#%d got %v, want %v", i, err, nil)
		}

		expectedPending := Events(events[i+1:])
		expectedState := Events(events[0 : i+1])

		// aggregate.pending should shrink with each sync event
		if !cmp.Equal(expectedPending, aggregate.pending, cmpopts.IgnoreUnexported(testAggregate{}, AggregateRoot{})) {
			t.Errorf("#%d got %v, want %v", i, aggregate.pending, expectedPending)
		}

		// aggregate.state should grow with each sync event
		if !cmp.Equal(expectedState, aggregate.state, cmpopts.IgnoreUnexported(testAggregate{}, AggregateRoot{})) {
			t.Errorf("#%d got %v, want %v", i, aggregate.state, expectedState)
		}
	}

	// should throw an error if there a no events to sync
	for i, v := range events {
		if err := (&aggregate).Sync(v); !errors.Is(err, ErrSyncNoEvents) {
			t.Fatalf("#%d got nil error, want %v", i, ErrSyncNoEvents)
		}
	}
}
