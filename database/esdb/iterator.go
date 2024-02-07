package esdb

import (
	"io"
	"strings"

	"github.com/EventStore/EventStore-Client-Go/esdb"
)

// Iterator represents an iterator allowing to iterate over stream of ledger.Events.
type Iterator struct {
	stream *esdb.ReadStream
	event  *esdb.ResolvedEvent
	err    error
}

func NewIterator(stream *esdb.ReadStream) *Iterator {
	i := &Iterator{
		stream: stream,
	}
	return i
}

// Close closes the stream
func (i *Iterator) Close() {
	if i.stream == nil {
		return
	}
	i.stream.Close()
}

// Next steps to the next event in the stream.
func (i *Iterator) Next() bool {
	if i.stream == nil {
		return false
	}

	eventESDB, err := i.stream.Recv()
	if err != nil {
		switch err {
		case io.EOF:
		default:
			i.err = err
		}
		return false
	}

	i.event = eventESDB

	return true
}

func (i *Iterator) Error() error {
	return i.err
}

// Value returns the event from the stream.
func (i *Iterator) Value() (*Event, error) {
	stream := strings.Split(i.event.Event.StreamID, streamSeparator) // TODO
	return &Event{
		AggregateID: stream[1],
		Version:     Version(i.event.Event.EventNumber),
		Type:        i.event.Event.EventType,
		Aggregate:   stream[0],
		Timestamp:   i.event.Event.CreatedDate,
		Data:        i.event.Event.Data,
		Metadata:    i.event.Event.UserMetadata,
	}, nil
}
