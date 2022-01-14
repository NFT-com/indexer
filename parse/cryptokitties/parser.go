package cryptokitties

import (
	"context"
	"log"
	"math/big"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/NFT-com/indexer/contracts/ethereum/cryptokitties"
	"github.com/NFT-com/indexer/parse"
)

func dispatch(nodeURL string) {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Println("failed to connect to node", err)
		os.Exit(1)
	}

	lambda.Start(Handler(client))
}

func Handler(client *ethclient.Client) func(ctx context.Context, nft parse.NFT) error {
	return func(ctx context.Context, nft parse.NFT) error {
		address := common.HexToAddress(nft.Address)
		instance, err := cryptokitties.NewMain(address, client)
		if err != nil {
			return err
		}

		metadata, err := instance.GetKitty(&bind.CallOpts{}, big.NewInt(nft.ID))
		if err != nil {
			return err
		}

		log.Println("Chain", nft.Chain, "Network", nft.Network, "Address", nft.Address, "NFT", nft.ID, "Genes", metadata.Genes)

		return nil
	}
}
