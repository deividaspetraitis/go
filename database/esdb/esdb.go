package esdb

import (
	"context"
	"fmt"
	"math"

	"github.com/deividaspetraitis/go/database"
	"github.com/deividaspetraitis/go/errors"

	"github.com/EventStore/EventStore-Client-Go/esdb"
)

const (
	count = math.MaxInt64 // maximum number of events
)

// DB is EventStore database client.
type Client struct {
	*esdb.Client
	contentType esdb.ContentType
}

func NewClient(config *database.Config) (*Client, error) {
	cfg, err := esdb.ParseConnectionString(dsn(config))
	if err != nil {
		return nil, err
	}

	client, err := esdb.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:      client,
		contentType: esdb.JsonContentType,
	}, nil
}

// dsn parses database.Config into DNS string.
func dsn(config *database.Config) string {
	return fmt.Sprintf("esdb://%s:%s@%s:%d?tls=false", config.Username, config.Password, config.Host, config.Port)
}

// Save stores given events into EventStore.
func (c *Client) Save(ctx context.Context, events []*Event) error {
	if len(events) == 0 {
		return nil
	}

	var data []esdb.EventData
	for _, v := range events {
		data = append(data, esdb.EventData{
			ContentType: esdb.JsonContentType,
			EventType:   v.Type,
			Data:        v.Data,
			Metadata:    v.Metadata,
		})
	}

	// last stored version is -1 than oldest event in the list
	version, err := parseFirstEventVersion(events)
	if err != nil {
		return err
	}

	// for the first event skip stream revision check
	var streamOptions esdb.AppendToStreamOptions
	if version > 1 {
		// EventStore events enumeration starts at 0, thus -2.
		streamOptions.ExpectedRevision = esdb.StreamRevision{Value: uint64(version) - 2}
	} else if version == 1 {
		streamOptions.ExpectedRevision = esdb.NoStream{}
	}

	if _, err := c.AppendToStream(context.Background(), stream(events), streamOptions, data...); err != nil {
		return err
	}

	return nil
}

// Get reads an stream of events for specific id and returns Iterator.
func (c *Client) Get(ctx context.Context, id string, aggregate string, afterVersion Version) (*Iterator, error) {
	stream, err := c.ReadStream(ctx, aggregate+"_"+id, esdb.ReadStreamOptions{
		From: esdb.StreamRevision{Value: uint64(afterVersion)},
	}, count)
	if err != nil {
		if errors.Is(err, esdb.ErrStreamNotFound) {
			return &Iterator{}, nil
		}
		return nil, err
	}
	return &Iterator{stream: stream}, nil
}

// parseFirstEventVersion parses and returns version from the first event in the list.
func parseFirstEventVersion(events []*Event) (Version, error) {
	if len(events) == 0 {
		return 0, errors.New("unable to determine version from empty events list")
	}
	return events[0].Version, nil
}
