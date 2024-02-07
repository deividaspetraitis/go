package database

import (
	"context"

	"github.com/deividaspetraitis/go/es"
)

// SaveAggregate stores aggregate into persistent storage.
// In the event of failure error will be returned.
type SaveAggregateFunc func(ctx context.Context, aggregate es.Aggregate) error

// GetAggregateFunc retrieves events from underlying database store and restores aggregate state.
// If Aggregate is not found an error will be returned.
type GetAggregateFunc[T any] func(ctx context.Context, aggregate es.Aggregate, id string) (T, error)
