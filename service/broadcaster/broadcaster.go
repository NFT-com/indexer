package broadcaster

const (
	HandlerKey = "handler"

	DiscoveryHandlerValue = "discovery"
	ParsingHandlerValue   = "parsing"
)

type Keys map[string]interface{}

func NewKeys() Keys {
	return map[string]interface{}{}
}

func (k Keys) WithHandler(value string) Keys {
	k[HandlerKey] = value

	return k
}

func (k Keys) HasHandler(value string) bool {
	rawHandler, ok := k[HandlerKey]
	if !ok {
		return false
	}

	handlerType, ok := rawHandler.(string)
	if !ok {
		return false
	}

	return handlerType == value
}
