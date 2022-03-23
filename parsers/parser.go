package parsers

type Parser interface {
	ParseRawEvent(rawEvent events.RawEvent) (events.Event, error)
}
