package pipeline

type Publisher interface {
	Publish(topic string, payload []byte) error
}

type BatchPublisher interface {
	MultiPublish(topic string, payloads [][]byte) error
}
