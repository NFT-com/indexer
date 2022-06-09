package ethereum

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/NFT-com/indexer/config/params"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

type SigningTransport struct {
	ctx         context.Context
	transport   http.RoundTripper
	credentials aws.Credentials
	region      string
}

func (s SigningTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	body, err := req.GetBody()
	if err != nil {
		return nil, fmt.Errorf("could not get request body: %w", err)
	}

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("could not read request body: %w", err)
	}

	hash := sha256.Sum256(data)
	payload := hex.EncodeToString(hash[:])

	signer := v4.NewSigner()
	err = signer.SignHTTP(s.ctx, s.credentials, req, payload, params.AWSManagedBlockchain, s.region, time.Now())
	if err != nil {
		return nil, fmt.Errorf("could not sign request: %w", err)
	}

	req.Header.Add("accept-encoding", "gzip, deflate")

	return s.transport.RoundTrip(req)
}
