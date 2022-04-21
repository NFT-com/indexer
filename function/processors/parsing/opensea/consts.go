package opensea

const (
	// This represents the type of this parser.
	openseaType = "0xc4109843e0b7d514e4c093114b863f8e7d8d9a458c372cd51bfe526b588006c9"

	defaultIndexDataLen = 3
	eventName           = "OrdersMatched"
	priceFieldName      = "price"

	// OrdersMatchedEventABI is extracted from the OpenSea marketplace contract ABI:
	// See https://etherscan.io/address/0x7f268357a8c2552623316e2562d90e642bb538e5#code
	eventABI = `[{"anonymous": false,"inputs": [{"indexed": false,"name": "buyHash","type": "bytes32"},{"indexed": false,"name": "sellHash","type": "bytes32"},{"indexed": true,"name": "maker","type": "address"},{"indexed": true,"name": "taker","type": "address"},{"indexed": false,"name": "price","type": "uint256"},{"indexed": true,"name": "metadata","type": "bytes32"}],"name": "OrdersMatched","type": "event"}]`
)
