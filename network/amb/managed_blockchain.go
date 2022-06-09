package amb

import (
	"fmt"
	"io"
	"net/http"
	"time"

	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

type RoundTripper struct {
	signer  *v4.Signer
	region  string
	service string
	wrap    http.RoundTripper
}

func NewRoundTripper(signer *v4.Signer, region string, service string, wrap http.RoundTripper) *RoundTripper {

	r := RoundTripper{
		signer:  signer,
		region:  region,
		service: service,
		wrap:    wrap,
	}

	return &r
}

func (t *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	reader, err := req.GetBody()
	if err != nil {
		return nil, fmt.Errorf("could not get request body: %w", err)
	}

	seeker, ok := reader.(io.ReadSeekCloser)
	if !ok {
		return nil, fmt.Errorf("could not cast request body (type: %T)", reader)
	}

	_, err = t.signer.Sign(req, seeker, t.service, t.region, time.Now())
	if err != nil {
		return nil, fmt.Errorf("could not sign request: %w", err)
	}

	req.Header.Add("accept-encoding", "gzip, deflate")

	return t.wrap.RoundTrip(req)
}
