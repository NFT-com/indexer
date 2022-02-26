package function

type Dispatcher interface {
	Dispatch(functionName string, payload []byte) error
}
