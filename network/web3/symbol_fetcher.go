package web3

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/models/abis"
)

type SymbolFetcher struct {
	client *ethclient.Client
}

func NewSymbolFetcher(client *ethclient.Client) *SymbolFetcher {

	u := SymbolFetcher{
		client: client,
	}

	return &u
}

func (u *SymbolFetcher) ERC20(ctx context.Context, address string) (string, error) {
	return u.fetch(ctx, address, nil, "symbol", abis.ERC20)
}

func (u *SymbolFetcher) fetch(ctx context.Context, address string, height *big.Int, name string, abi abi.ABI) (string, error) {

	input, err := abi.Pack(name)
	if err != nil {
		return "", fmt.Errorf("could not pack input: %w", err)
	}

	ethAddress := common.HexToAddress(address)
	msg := ethereum.CallMsg{From: common.Address{}, To: &ethAddress, Data: input}
	output, err := u.client.CallContract(ctx, msg, height)
	if err != nil {
		return "", fmt.Errorf("could not execute contract call: %w", err)
	}

	fields, err := abi.Unpack(name, output)
	if err != nil {
		return "", fmt.Errorf("could not unpack output: %w", err)
	}

	if len(fields) != 1 {
		return "", fmt.Errorf("invalid number of fields (have: %d, want: 1)", len(fields))
	}

	symbol, ok := fields[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid field type (have: %T, want: string)", fields[0])
	}

	return symbol, nil
}
