package aws

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws/signer/v4"
)

type InjectorRoundTripper struct {
	signer       *v4.Signer
	region       string
	service      string
	roundTripper http.RoundTripper
}

func NewInjectorRoundTripper(signer *v4.Signer, region string, service string, roundTripper http.RoundTripper) InjectorRoundTripper {
	i := InjectorRoundTripper{
		signer:       signer,
		region:       region,
		service:      service,
		roundTripper: roundTripper,
	}

	return i
}

func (t InjectorRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	reader, err := req.GetBody()
	if err != nil {
		return nil, fmt.Errorf("could not get the request body: %w", err)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read the request body: %w", err)
	}

	_, err = t.signer.Sign(req, bytes.NewReader(body), t.service, t.region, time.Now())
	if err != nil {
		return nil, fmt.Errorf("could not sign request: %w", err)
	}

	req.Header.Add("accept-encoding", "gzip, deflate")
	return t.roundTripper.RoundTrip(req)
}
