package addition

const (
	// represents what type of the processor action this is
	processorType = "addition"
	// represents the type of standard the processor is able to parse
	processorStandard = "ERC721"

	dateDisplayType = "date"

	// callSender represents the address that will be set as signer of the get request to the node.
	callSender           = "0xd45FCC235228431812C615F1D4Be4914b6D37593"
	tokenURIFunctionName = "tokenURI"
	uriFunctionABI       = `[
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "tokenId",
				"type": "uint256"
			}
		],
		"name": "tokenURI",
		"outputs": [
			{
				"internalType": "string",
				"name": "",
				"type": "string"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`
)
