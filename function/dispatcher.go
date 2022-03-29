package function

type Invoker interface {
	Invoke(name string, payload []byte) ([]byte, error)
}
