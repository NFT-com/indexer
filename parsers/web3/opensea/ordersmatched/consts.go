package main

const (
	eventName      = "OrdersMatched"
	priceFieldName = "price"

	// OrdersMatchedEventABI is extracted from the OpenSea marketplace contract ABI:
	// See https://etherscan.io/address/0x7f268357a8c2552623316e2562d90e642bb538e5#code
	eventABI = `[{"anonymous": false,"inputs": [{"indexed": false,"name": "buyHash","type": "bytes32"},{"indexed": false,"name": "sellHash","type": "bytes32"},{"indexed": true,"name": "maker","type": "address"},{"indexed": true,"name": "taker","type": "address"},{"indexed": false,"name": "price","type": "uint256"},{"indexed": true,"name": "metadata","type": "bytes32"}],"name": "OrdersMatched","type": "event"}]`
)
