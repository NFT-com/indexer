package batchtransfer

const (
	// This represents the type of this parser.
	transferType = "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb"

	// The number of indexed fields required in the array.
	defaultIndexDataLen = 3

	// This represents the null address where nfts mint from or burn to.
	zeroValueAddress = "0x0000000000000000000000000000000000000000"

	// Event name in the ABI.
	eventName = "TransferBatch"
	// IDs field name in the event arguments.
	idsFieldName = "_ids"
	// Values field name in the event arguments.
	valuesFieldName = "_values"
	// eventABI is extracted from the OpenSea marketplace contract ABI:
	// See https://etherscan.io/address/0xc36cf0cfcb5d905b8b513860db0cfe63f6cf9f5c#code.
	eventABI = `[{"anonymous": false,"inputs": [{"indexed": true,"internalType": "address","name": "_operator","type": "address"},{"indexed": true,"internalType": "address","name": "_from","type": "address"},{"indexed": true,"internalType": "address","name": "_to","type": "address"},{"indexed": false,"internalType": "uint256[]","name": "_ids","type": "uint256[]"},{"indexed": false,"internalType": "uint256[]","name": "_values","type": "uint256[]"}],"name": "TransferBatch","type": "event"}]`
)
