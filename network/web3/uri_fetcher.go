package web3

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

type URIFetcher struct {
	client ethclient.Client
}

func NewURIFetcher(client ethclient.Client) *URIFetcher {

	u := URIFetcher{
		client: client,
	}

	return &u
}

func (u *URIFetcher) ERC721(address string) (string, error) {
}

func (u *URIFetcher) ERC1155(address string) (string, error) {
}
