package lambdas

type URIFetcher interface {
	ERC721(address string, tokenID string) (string, error)
	ERC1155(address string, tokenID string) (string, error)
}
