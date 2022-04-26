package uri

const (
	// This represents the type of this parser.
	uriType = "0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b"

	// The number of indexed fields required in the array
	defaultIndexDataLen = 1

	// event name in the ABI
	eventName = "URI"
	// uri field name in the event arguments
	uriFieldName = "value"

	// eventABI is extracted from the OpenSea marketplace contract ABI:
	// See https://etherscan.io/address/0xc36cf0cfcb5d905b8b513860db0cfe63f6cf9f5c#code
	eventABI = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"string","name":"value","type":"string"},{"indexed":true,"internalType":"uint256","name":"id","type":"uint256"}],"name":"URI","type":"event"}]`
)
