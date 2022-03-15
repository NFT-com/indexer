package function

type Dispatcher interface {
	Dispatch(name string, payload []byte) (interface{}, error)
}
