package esdb

import "time"

// Version represents event version
type Version uint64

// Event represents Event entity.
type Event struct {
	AggregateID string
	Version     Version
	Aggregate   string
	Type        string
	Timestamp   time.Time
	Data        []byte
	Metadata    []byte
}
