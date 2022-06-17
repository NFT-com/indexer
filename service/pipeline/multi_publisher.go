package pipeline

type MultiPublisher interface {
	MultiPublish(topic string, payloads [][]byte) error
}
