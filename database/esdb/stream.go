package esdb

var streamSeparator = "_"

func stream(events []*Event) string {
	return events[0].Aggregate + "_" + events[0].AggregateID // TODO
}
