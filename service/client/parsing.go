package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/request"
)

func (c *Client) SubscribeNewParsingJob(subscriberType string, parsingJobs chan []jobs.Parsing) error {
	params := url.Values{}
	if subscriberType != SubscriberTypeAllJobs {
		params.Set("status", subscriberType)
	}

	url := c.config.jobsWebsocket
	url.Path = path.Join("ws", parsingBasePath)
	url.RawQuery = params.Encode()

	connection, _, err := c.config.dialer.Dial(url.String(), nil)
	if err != nil {
		return fmt.Errorf("could not dial websocket: %w", err)
	}

	internalClose := make(chan struct{})
	connection.SetCloseHandler(func(code int, text string) error {
		c.log.Info().Int("code", code).Str("text", text).Msg("parsing jobs websocket connection closed")
		close(internalClose)
		return nil
	})

	go func() {
		for {
			select {
			case <-c.close:
				return
			case <-internalClose:
				return
			default:
			}

			var jobs []jobs.Parsing
			err := connection.ReadJSON(&jobs)
			if err != nil {
				c.log.Error().Err(err).Msg("could not read message socket")
				continue
			}

			parsingJobs <- jobs
		}
	}()

	return nil
}

func (c *Client) CreateParsingJob(job jobs.Parsing) (*jobs.Parsing, error) {
	req := request.Parsing{
		ChainURL:     job.ChainURL,
		ChainID:      job.ChainID,
		ChainType:    job.ChainType,
		BlockNumber:  job.BlockNumber,
		Address:      job.Address,
		StandardType: job.StandardType,
		EventType:    job.EventType,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %w", err)
	}

	url := c.config.jobsAPI
	url.Path = parsingBasePath

	resp, err := c.config.client.Post(url.String(), jsonContentType, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("could not create job: got status code %v", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close response body: %w", err)
	}

	newJob := jobs.Parsing{}
	err = json.Unmarshal(responseBody, &newJob)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &newJob, nil
}

func (c *Client) CreateParsingJobs(jobList []jobs.Parsing) error {
	requestJobs := make([]request.Parsing, 0, len(jobList))
	for _, job := range jobList {
		requestJob := request.Parsing{
			ChainURL:     job.ChainURL,
			ChainID:      job.ChainID,
			ChainType:    job.ChainType,
			BlockNumber:  job.BlockNumber,
			Address:      job.Address,
			StandardType: job.StandardType,
		}

		requestJobs = append(requestJobs, requestJob)
	}

	req := request.Parsings{
		Jobs: requestJobs,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	url := c.config.jobsAPI
	url.Path = path.Join(parsingBasePath, "batch")

	_, err = c.config.client.Post(url.String(), jsonContentType, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}

	return nil
}

func (c *Client) ListParsingJobs(status jobs.Status) ([]jobs.Parsing, error) {
	params := url.Values{}
	params.Set("status", string(status))

	url := c.config.jobsAPI
	url.Path = parsingBasePath
	url.RawQuery = params.Encode()

	resp, err := c.config.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("could not list job: got status code %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close response body: %w", err)
	}

	jobList := make([]jobs.Parsing, 0)
	err = json.Unmarshal(body, &jobList)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return jobList, nil
}

func (c *Client) GetParsingJob(id string) (*jobs.Parsing, error) {
	url := c.config.jobsAPI
	url.Path = path.Join(parsingBasePath, id)

	resp, err := c.config.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("could not get job: got status code %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close response body: %w", err)
	}

	job := jobs.Parsing{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &job, nil
}

func (c *Client) GetHighestBlockNumberParsingJob(chainURL, chainType, address, standardType, eventType string) (*jobs.Parsing, error) {
	params := url.Values{}
	params.Set("chain_url", chainURL)
	params.Set("chain_type", chainType)
	params.Set("address", address)
	params.Set("standard_type", standardType)
	params.Set("event_type", eventType)

	url := c.config.jobsAPI
	url.Path = path.Join(parsingBasePath, "highest")
	url.RawQuery = params.Encode()

	resp, err := c.config.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("could not highest block number job: got status code %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close response body: %w", err)
	}

	job := jobs.Parsing{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &job, nil
}

func (c *Client) UpdateParsingJobStatus(id string, status jobs.Status) error {
	requestBody := request.Status{
		Status: string(status),
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	url := c.config.jobsAPI
	url.Path = path.Join(parsingBasePath, id)

	req, err := http.NewRequest(http.MethodPatch, url.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Add(contentTypeHeaderName, jsonContentType)

	resp, err := c.config.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("could not update job: got status code %v", resp.StatusCode)
	}

	return nil
}
