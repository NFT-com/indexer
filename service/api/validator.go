package api

type Validator interface {
	Request(interface{}) error
}
