package creator

type Publisher interface {
	Publish(topic string, payload []byte) error
}
