package processors

const (
	// TODO: remove whitespaces

	ERC1155TransferABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"operator","type":"address"},{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"id","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"TransferSingle","type":"event"}]`
	ERC1155BatchABI    = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"_operator","type":"address"},{"indexed":true,"internalType":"address","name":"_from","type":"address"},{"indexed":true,"internalType":"address","name":"_to","type":"address"},{"indexed":false,"internalType":"uint256[]","name":"_ids","type":"uint256[]"},{"indexed":false,"internalType": "uint256[]","name": "_values","type": "uint256[]"}],"name": "TransferBatch","type": "event"}]`

	// https://etherscan.io/address/0xc36cf0cfcb5d905b8b513860db0cfe63f6cf9f5c#code
	OpenSeaTradeABI = `[{"anonymous": false,"inputs": [{"indexed": false,"name": "buyHash","type": "bytes32"},{"indexed": false,"name": "sellHash","type": "bytes32"},{"indexed": true,"name": "maker","type": "address"},{"indexed": true,"name": "taker","type": "address"},{"indexed": false,"name": "price","type": "uint256"},{"indexed": true,"name": "metadata","type": "bytes32"}],"name": "OrdersMatched","type": "event"}]`
)
