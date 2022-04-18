package broadcaster

const (
	handlerKey = "handler"
	statusKey  = "status"

	// DiscoveryHandlerValue represents the discovery handler value.
	DiscoveryHandlerValue = "discovery"
	// ParsingHandlerValue represents the parsing handler value.
	ParsingHandlerValue = "parsing"
	// AdditionHandlerValue represents the addition handler value.
	AdditionHandlerValue = "addition"
	// CreateStatusValue represents the create status value.
	CreateStatusValue = "create"
	// UpdateStatusValue represents the update status value.
	UpdateStatusValue = "update"
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

// HasHandlerKey returns if the keys have the key set with value.
func HasHandlerKey(keys map[string]interface{}) bool {
	_, ok := keys[handlerKey]
	if !ok {
		return false
	}

	return true
}

// WithStatus returns new keys with the status key
func WithStatus(keys map[string]interface{}, value string) map[string]interface{} {
	keys[statusKey] = value

	return keys
}

// HasStatus returns if the keys have the status key value set.
func HasStatus(keys map[string]interface{}, value string) bool {
	rawStatus, ok := keys[statusKey]
	if !ok {
		return false
	}

	statusType, ok := rawStatus.(string)
	if !ok {
		return false
	}

	return statusType == value
}

// HasStatusKey returns if the keys have the status key set with value.
func HasStatusKey(keys map[string]interface{}) bool {
	_, ok := keys[statusKey]
	if !ok {
		return false
	}

	return true
}
