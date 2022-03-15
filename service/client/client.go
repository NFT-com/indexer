package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/job"
	"github.com/NFT-com/indexer/service/request"
)

const (
	DiscoveryBasePath = "discoveries"
	ParsingBasePath   = "parsings"

	ContentTypeHeaderName = "content-type"
	JsonContentType       = "application/json"
)

type Client struct {
	log        zerolog.Logger
	wsClient   *websocket.Dialer
	httpClient *http.Client
	options    *options
	close      chan struct{}
}

func NewClient(log zerolog.Logger, optionList OptionsList) *Client {
	opts := defaultOptions()
	optionList.Apply(opts)

	c := Client{
		log:        log.With().Str("component", "api_client").Logger(),
		wsClient:   opts.wsDialer,
		httpClient: opts.httpClient,
		options:    opts,
		close:      make(chan struct{}),
	}

	return &c
}

func (c *Client) SubscribeNewDiscoveryJob(discoveryJobs chan job.Discovery) error {
	requestURL := fmt.Sprintf("%s/ws/%s", c.options.websocketURL.String(), DiscoveryBasePath)
	connection, _, err := c.wsClient.Dial(requestURL, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-c.close:
				return
			default:
				newDiscoveryJob := job.Discovery{}
				err := connection.ReadJSON(&newDiscoveryJob)
				if err != nil {
					c.log.Error().Err(err).Msg("failed to read message socket")
					return
				}

				discoveryJobs <- newDiscoveryJob
			}
		}
	}()

	return nil
}

func (c *Client) SubscribeNewParsingJob(parsingJobs chan job.Parsing) error {
	requestURL := fmt.Sprintf("%s/ws/%s", c.options.websocketURL.String(), ParsingBasePath)
	connection, _, err := c.wsClient.Dial(requestURL, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-c.close:
				return
			default:
				newParsingJob := job.Parsing{}
				err := connection.ReadJSON(&newParsingJob)
				if err != nil {
					c.log.Error().Err(err).Msg("failed to read message socket")
					return
				}

				parsingJobs <- newParsingJob
			}
		}
	}()

	return nil
}

func (c *Client) CreateDiscoveryJob(discoveryJob job.Discovery) (job.Discovery, error) {
	body, err := json.Marshal(discoveryJob)
	if err != nil {
		return job.Discovery{}, err
	}

	requestURL := fmt.Sprintf("%s/%s", c.options.httpURL.String(), DiscoveryBasePath)
	resp, err := c.httpClient.Post(requestURL, JsonContentType, bytes.NewReader(body))
	if err != nil {
		return job.Discovery{}, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return job.Discovery{}, err
	}

	err = resp.Body.Close()
	if err != nil {
		return job.Discovery{}, err
	}

	newJob := job.Discovery{}
	err = json.Unmarshal(responseBody, &newJob)
	if err != nil {
		return job.Discovery{}, err
	}

	return newJob, nil
}

func (c *Client) ListDiscoveryJobs(status job.Status) ([]job.Discovery, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s?status=%s", c.options.httpURL.String(), DiscoveryBasePath, status))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	jobList := make([]job.Discovery, 0)
	err = json.Unmarshal(body, &jobList)
	if err != nil {
		return nil, err
	}

	return jobList, nil
}

func (c *Client) GetDiscoveryJob(jobID string) (job.Discovery, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), DiscoveryBasePath, jobID))
	if err != nil {
		return job.Discovery{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return job.Discovery{}, err
	}

	err = resp.Body.Close()
	if err != nil {
		return job.Discovery{}, err
	}

	discoveryJob := job.Discovery{}
	err = json.Unmarshal(body, &discoveryJob)
	if err != nil {
		return job.Discovery{}, err
	}

	return discoveryJob, nil
}

func (c *Client) UpdateDiscoveryJobState(jobID string, jobStatus job.Status) error {
	requestBody := request.Status{Status: string(jobStatus)}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), DiscoveryBasePath, jobID)
	req, err := http.NewRequest(http.MethodPatch, requestURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add(ContentTypeHeaderName, JsonContentType)

	_, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RequeueDiscoveryJob(jobID string) (job.Discovery, error) {
	resp, err := c.httpClient.Post(fmt.Sprintf("%s/%s/%s/requeue", c.options.httpURL.String(), DiscoveryBasePath, jobID), JsonContentType, nil)
	if err != nil {
		return job.Discovery{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return job.Discovery{}, err
	}

	err = resp.Body.Close()
	if err != nil {
		return job.Discovery{}, err
	}

	discoveryJob := job.Discovery{}
	err = json.Unmarshal(body, &discoveryJob)
	if err != nil {
		return job.Discovery{}, err
	}

	return discoveryJob, nil
}

func (c *Client) CreateParsingJob(parsingJob job.Parsing) (job.Parsing, error) {
	body, err := json.Marshal(parsingJob)
	if err != nil {
		return job.Parsing{}, err
	}

	requestURL := fmt.Sprintf("%s/%s", c.options.httpURL.String(), ParsingBasePath)
	resp, err := c.httpClient.Post(requestURL, JsonContentType, bytes.NewReader(body))
	if err != nil {
		return job.Parsing{}, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return job.Parsing{}, err
	}

	err = resp.Body.Close()
	if err != nil {
		return job.Parsing{}, err
	}

	newJob := job.Parsing{}
	err = json.Unmarshal(responseBody, &newJob)
	if err != nil {
		return job.Parsing{}, err
	}

	return newJob, nil
}

func (c *Client) ListParsingJobs(status job.Status) ([]job.Parsing, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s?status=%s", c.options.httpURL.String(), ParsingBasePath, status))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	jobList := make([]job.Parsing, 0)
	err = json.Unmarshal(body, &jobList)
	if err != nil {
		return nil, err
	}

	return jobList, nil
}

func (c *Client) GetParsingJob(jobID string) (job.Parsing, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), ParsingBasePath, jobID))
	if err != nil {
		return job.Parsing{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return job.Parsing{}, err
	}

	err = resp.Body.Close()
	if err != nil {
		return job.Parsing{}, err
	}

	parsingJob := job.Parsing{}
	err = json.Unmarshal(body, &parsingJob)
	if err != nil {
		return job.Parsing{}, err
	}

	return parsingJob, nil
}

func (c *Client) UpdateParsingJobState(jobID string, jobStatus job.Status) error {
	requestBody := request.Status{Status: string(jobStatus)}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), ParsingBasePath, jobID)
	req, err := http.NewRequest(http.MethodPatch, requestURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Add(ContentTypeHeaderName, JsonContentType)

	_, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RequeueParsingJob(jobID string) (job.Parsing, error) {
	resp, err := c.httpClient.Post(fmt.Sprintf("%s/%s/%s/requeue", c.options.httpURL.String(), ParsingBasePath, jobID), JsonContentType, nil)
	if err != nil {
		return job.Parsing{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return job.Parsing{}, err
	}

	err = resp.Body.Close()
	if err != nil {
		return job.Parsing{}, err
	}

	parsingJob := job.Parsing{}
	err = json.Unmarshal(body, &parsingJob)
	if err != nil {
		return job.Parsing{}, err
	}

	return parsingJob, nil
}

func (c *Client) Close() {
	close(c.close)
}
