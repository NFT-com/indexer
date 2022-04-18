package main

import (
	"crypto/sha256"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		h := sha256.New()
		h.Write([]byte("web3-erc721-0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"))
		web3Erc721TransferName := fmt.Sprintf("%x", h.Sum(nil))
		_, err := lambda.NewFunction(ctx, web3Erc721TransferName, &lambda.FunctionArgs{
			Code:        pulumi.NewFileArchive("erc721_transfer.zip"),
			Name:        pulumi.String(web3Erc721TransferName),
			Handler:     pulumi.String("worker"),
			Description: pulumi.String("web3-erc721-0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			Role:        pulumi.String("arn:aws:iam::567721414829:role/lamda-role"),
			Runtime:     pulumi.String("go1.x"),
		})
		if err != nil {
			return err
		}

		h.Reset()
		h.Write([]byte("web3-opensea-0xc4109843e0b7d514e4c093114b863f8e7d8d9a458c372cd51bfe526b588006c9"))
		web3OpenSeaOrdersMatchedName := fmt.Sprintf("%x", h.Sum(nil))
		_, err = lambda.NewFunction(ctx, web3OpenSeaOrdersMatchedName, &lambda.FunctionArgs{
			Code:        pulumi.NewFileArchive("opensea_ordersmatched.zip"),
			Name:        pulumi.String(web3OpenSeaOrdersMatchedName),
			Handler:     pulumi.String("worker"),
			Description: pulumi.String("web3-opensea-0xc4109843e0b7d514e4c093114b863f8e7d8d9a458c372cd51bfe526b588006c9"),
			Role:        pulumi.String("arn:aws:iam::567721414829:role/lamda-role"),
			Runtime:     pulumi.String("go1.x"),
		})
		if err != nil {
			return err
		}

		h.Reset()
		h.Write([]byte("web3-erc721"))
		web3ERC721Addition := fmt.Sprintf("%x", h.Sum(nil))
		_, err = lambda.NewFunction(ctx, web3ERC721Addition, &lambda.FunctionArgs{
			Code:        pulumi.NewFileArchive("addition.zip"),
			Name:        pulumi.String(web3ERC721Addition),
			Handler:     pulumi.String("worker"),
			Description: pulumi.String("web3-ERC721-addition"),
			Role:        pulumi.String("arn:aws:iam::567721414829:role/lamda-role"),
			Runtime:     pulumi.String("go1.x"),
		})
		if err != nil {
			return err
		}

		return nil
	})
}
