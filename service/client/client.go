package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/jobs"
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

func (c *Client) SubscribeNewDiscoveryJob(discoveryJobs chan jobs.Discovery) error {
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
				newDiscoveryJob := jobs.Discovery{}
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

func (c *Client) SubscribeNewParsingJob(parsingJobs chan jobs.Parsing) error {
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
				newParsingJob := jobs.Parsing{}
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

func (c *Client) CreateDiscoveryJob(job jobs.Discovery) (*jobs.Discovery, error) {
	req := request.Discovery{
		ChainURL:     job.ChainURL,
		ChainType:    job.ChainType,
		BlockNumber:  job.BlockNumber,
		Addresses:    job.Addresses,
		StandardType: job.StandardType,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	requestURL := fmt.Sprintf("%s/%s", c.options.httpURL.String(), DiscoveryBasePath)
	resp, err := c.httpClient.Post(requestURL, JsonContentType, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	newJob := jobs.Discovery{}
	err = json.Unmarshal(responseBody, &newJob)
	if err != nil {
		return nil, err
	}

	return &newJob, nil
}

func (c *Client) ListDiscoveryJobs(status jobs.Status) ([]jobs.Discovery, error) {
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

	jobList := make([]jobs.Discovery, 0)
	err = json.Unmarshal(body, &jobList)
	if err != nil {
		return nil, err
	}

	return jobList, nil
}

func (c *Client) GetDiscoveryJob(id string) (*jobs.Discovery, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), DiscoveryBasePath, id))
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

	job := jobs.Discovery{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *Client) UpdateDiscoveryJobState(id string, status jobs.Status) error {
	requestBody := request.Status{Status: string(status)}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), DiscoveryBasePath, id)
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

func (c *Client) RequeueDiscoveryJob(id string) (*jobs.Discovery, error) {
	resp, err := c.httpClient.Post(fmt.Sprintf("%s/%s/%s/requeue", c.options.httpURL.String(), DiscoveryBasePath, id), JsonContentType, nil)
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

	job := jobs.Discovery{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *Client) CreateParsingJob(job jobs.Parsing) (*jobs.Parsing, error) {
	req := request.Parsing{
		ChainURL:     job.ChainURL,
		ChainType:    job.ChainType,
		BlockNumber:  job.BlockNumber,
		Address:      job.Address,
		StandardType: job.StandardType,
		EventType:    job.EventType,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	requestURL := fmt.Sprintf("%s/%s", c.options.httpURL.String(), ParsingBasePath)
	resp, err := c.httpClient.Post(requestURL, JsonContentType, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	newJob := jobs.Parsing{}
	err = json.Unmarshal(responseBody, &newJob)
	if err != nil {
		return nil, err
	}

	return &newJob, nil
}

func (c *Client) ListParsingJobs(status jobs.Status) ([]jobs.Parsing, error) {
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

	jobList := make([]jobs.Parsing, 0)
	err = json.Unmarshal(body, &jobList)
	if err != nil {
		return nil, err
	}

	return jobList, nil
}

func (c *Client) GetParsingJob(id string) (*jobs.Parsing, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), ParsingBasePath, id))
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

	job := jobs.Parsing{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *Client) UpdateParsingJobState(id string, status jobs.Status) error {
	requestBody := request.Status{Status: string(status)}
	body, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), ParsingBasePath, id)
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

func (c *Client) RequeueParsingJob(id string) (*jobs.Parsing, error) {
	resp, err := c.httpClient.Post(fmt.Sprintf("%s/%s/%s/requeue", c.options.httpURL.String(), ParsingBasePath, id), JsonContentType, nil)
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

	job := jobs.Parsing{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (c *Client) Close() {
	close(c.close)
}
