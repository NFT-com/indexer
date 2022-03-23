package api

// Validator represent a request validator interface.
type Validator interface {
	Request(interface{}) error
}
