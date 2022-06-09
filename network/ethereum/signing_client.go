package ethereum

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func NewSigningClient(ctx context.Context, url string, cfg aws.Config) (*ethclient.Client, error) {

	credentials, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve AWS credentials: %w", err)
	}

	client := http.Client{
		Transport: SigningTransport{
			ctx:         ctx,
			credentials: credentials,
			region:      cfg.Region,
		},
	}

	rpc, err := rpc.DialHTTPWithClient(url, &client)
	if err != nil {
		return nil, fmt.Errorf("could not connect to JSON RPC API: %w", err)
	}

	api := ethclient.NewClient(rpc)

	return api, nil
}
