package broadcaster

const (
	handlerKey = "handler"

	// DiscoveryHandlerValue represents the discovery handler value.
	DiscoveryHandlerValue = "discovery"
	// ParsingHandlerValue represents the parsing handler value.
	ParsingHandlerValue = "parsing"
)

// WithHandler returns new keys with the handler key
func WithHandler(keys map[string]interface{}, value string) map[string]interface{} {
	keys[handlerKey] = value

	return keys
}

// HasHandler returns if the keys have the handler key value set.
func HasHandler(keys map[string]interface{}, value string) bool {
	rawHandler, ok := keys[handlerKey]
	if !ok {
		return false
	}

	handlerType, ok := rawHandler.(string)
	if !ok {
		return false
	}

	return handlerType == value
}
