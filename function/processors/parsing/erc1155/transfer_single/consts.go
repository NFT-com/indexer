package singletransfer

const (
	// This represents the type of this parser.
	transferType = "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"

	// The number of indexed fields required in the array.
	defaultIndexDataLen = 3

	// This represents the null address where nfts mint from or burn to.
	zeroValueAddress = "0x0000000000000000000000000000000000000000"

	// Event name in the ABI.
	eventName = "TransferSingle"
	// ID field name in the event arguments.
	idFieldName = "id"
	// Value field name in the event arguments.
	valueFieldName = "value"
	// eventABI is extracted from the OpenSea marketplace contract ABI:
	// See https://etherscan.io/address/0xc36cf0cfcb5d905b8b513860db0cfe63f6cf9f5c#code.
	eventABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"operator","type":"address"},{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"TransferSingle","type":"event"}]`
)
