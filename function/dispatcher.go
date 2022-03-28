package function

type Dispatcher interface {
	Invoke(name string, payload []byte) ([]byte, error)
}
