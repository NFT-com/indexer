package pipeline

type Publisher interface {
	Publish(topic string, payload []byte) error
	MultiPublish(topic string, payloads [][]byte) error
}
