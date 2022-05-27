package web3

import (
	"context"
	"fmt"
	"math/big"
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

func (u *URIFetcher) ERC721(ctx context.Context, address string, height uint64, tokenID string) (string, error) {
	return u.fetch(ctx, address, tokenID, height, "tokenURI", abis.ERC721)
}

func (u *URIFetcher) ERC1155(ctx context.Context, address string, height uint64, tokenID string) (string, error) {
	return u.fetch(ctx, address, tokenID, height, "uri", abis.ERC1155)
}

func (u *URIFetcher) fetch(ctx context.Context, address string, tokenID string, height uint64, name string, abi abi.ABI) (string, error) {

	id, ok := big.NewInt(0).SetString(tokenID, 10)
	if !ok {
		return "", fmt.Errorf("could not convert token ID to integer")
	}

	input, err := abi.Pack(name, id)
	if err != nil {
		return "", fmt.Errorf("could not pack input: %w", err)
	}

	ethAddress := common.HexToAddress(address)
	msg := ethereum.CallMsg{From: common.Address{}, To: &ethAddress, Data: input}
	output, err := u.client.CallContract(ctx, msg, nil)
	if err != nil && !strings.Contains(err.Error(), "nonexistent token") {
		return "", fmt.Errorf("could not call present contract: %w", err)
	}

	if strings.Contains(err.Error(), "nonexistent token") {
		output, err = u.client.CallContract(ctx, msg, big.NewInt(0).SetUint64(height))
		if err != nil {
			return "", fmt.Errorf("could not call past contract: %w", err)
		}
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
