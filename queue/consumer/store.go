package consumer

type Store interface {
	InsertHistory(event events.Event) error
}
