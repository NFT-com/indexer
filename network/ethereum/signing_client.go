package ethereum

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func NewSigningClient(ctx context.Context, url string, cfg aws.Config) (*ethclient.Client, func(), error) {

	credentials, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("could not retrieve AWS credentials: %w", err)
	}

	dial := net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: -1,
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dial.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          1,
		IdleConnTimeout:       5 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := http.Client{
		Transport: SigningTransport{
			ctx:         ctx,
			transport:   transport,
			credentials: credentials,
			region:      cfg.Region,
		},
	}

	rpc, err := rpc.DialHTTPWithClient(url, &client)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to JSON RPC API: %w", err)
	}

	api := ethclient.NewClient(rpc)

	return api, client.CloseIdleConnections, nil
}
