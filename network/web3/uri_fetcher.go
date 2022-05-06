package web3

import (
	"context"
	"fmt"
	"strings"

	"github.com/NFT-com/indexer/models/abis"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type URIFetcher struct {
	client *ethclient.Client
}

func NewURIFetcher(client *ethclient.Client) *URIFetcher {

	u := URIFetcher{
		client: client,
	}

	return &u
}

func (u *URIFetcher) ERC721(ctx context.Context, address string, tokenID string) (string, error) {
	return u.fetch(ctx, address, tokenID, "tokenURI", abis.ERC721)
}

func (u *URIFetcher) ERC1155(ctx context.Context, address string, tokenID string) (string, error) {
	return u.fetch(ctx, address, tokenID, "uri", abis.ERC1155)
}

func (u *URIFetcher) fetch(ctx context.Context, address string, tokenID string, name string, abi abi.ABI) (string, error) {

	input, err := abi.Pack(name, tokenID)
	if err != nil {
		return "", fmt.Errorf("could not pack input: %w", err)
	}

	ethAddress := common.HexToAddress(address)
	msg := ethereum.CallMsg{From: common.Address{}, To: &ethAddress, Data: input}
	output, err := u.client.CallContract(ctx, msg, nil)
	if err != nil {
		return "", fmt.Errorf("could not call contract: %w", err)
	}

	fields, err := abi.Unpack(name, output)
	if err != nil {
		return "", fmt.Errorf("could not unpack output: %w", err)
	}

	if len(fields) != 1 {
		return "", fmt.Errorf("invalid number of fields (have: %d, want: 1)", len(fields))
	}

	uri, ok := fields[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid field type (have: %T, want: string)", fields[0])
	}

	uri = strings.ReplaceAll(uri, "ipfs://", "https://ipfs.io/ipfs/")

	return uri, nil
}