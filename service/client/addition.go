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

func (c *Client) SubscribeNewAdditionJob(subscriberType string, additionJobs chan []jobs.Addition) error {
	params := url.Values{}
	if subscriberType != SubscriberTypeAllJobs {
		params.Set("status", subscriberType)
	}

	url := c.config.jobsWebsocket
	url.Path = path.Join("ws", additionBasePath)
	url.RawQuery = params.Encode()

	connection, _, err := c.config.dialer.Dial(url.String(), nil)
	if err != nil {
		return fmt.Errorf("could not dial websocket: %w", err)
	}

	internalClose := make(chan struct{})
	connection.SetCloseHandler(func(code int, text string) error {
		c.log.Info().Int("code", code).Str("text", text).Msg("addition jobs websocket connection closed")
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

			var jobs []jobs.Addition
			err := connection.ReadJSON(&jobs)
			if err != nil {
				c.log.Error().Err(err).Msg("could not read message socket")
				continue
			}

			additionJobs <- jobs
		}
	}()

	return nil
}

func (c *Client) CreateAdditionJob(job jobs.Addition) (*jobs.Addition, error) {
	req := request.Addition{
		ChainURL:     job.ChainURL,
		ChainID:      job.ChainID,
		ChainType:    job.ChainType,
		BlockNumber:  job.BlockNumber,
		Address:      job.Address,
		StandardType: job.StandardType,
		TokenID:      job.TokenID,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %w", err)
	}

	url := c.config.jobsAPI
	url.Path = additionBasePath

	resp, err := c.config.client.Post(url.String(), jsonContentType, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("could not create job: got status code %d", resp.StatusCode)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}
	defer resp.Body.Close()

	var newJob jobs.Addition
	err = json.Unmarshal(responseBody, &newJob)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &newJob, nil
}

func (c *Client) CreateAdditionJobs(jobList []jobs.Addition) error {
	requestJobs := make([]request.Addition, 0, len(jobList))
	for _, job := range jobList {
		requestJob := request.Addition{
			ChainURL:     job.ChainURL,
			ChainID:      job.ChainID,
			ChainType:    job.ChainType,
			BlockNumber:  job.BlockNumber,
			Address:      job.Address,
			StandardType: job.StandardType,
			TokenID:      job.TokenID,
		}

		requestJobs = append(requestJobs, requestJob)
	}

	req := request.Additions{
		Jobs: requestJobs,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	url := c.config.jobsAPI
	url.Path = path.Join(additionBasePath, "batch")

	_, err = c.config.client.Post(url.String(), jsonContentType, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}

	return nil
}

func (c *Client) ListAdditionJobs(status jobs.Status) ([]jobs.Addition, error) {
	params := url.Values{}
	params.Set("status", string(status))

	url := c.config.jobsAPI
	url.Path = additionBasePath
	url.RawQuery = params.Encode()

	resp, err := c.config.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("could not list job: got status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}
	defer resp.Body.Close()

	jobList := make([]jobs.Addition, 0)
	err = json.Unmarshal(body, &jobList)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return jobList, nil
}

func (c *Client) GetAdditionJob(id string) (*jobs.Addition, error) {
	url := c.config.jobsAPI
	url.Path = path.Join(additionBasePath, id)

	resp, err := c.config.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("could not get job: got status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}
	defer resp.Body.Close()

	var job jobs.Addition
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &job, nil
}

func (c *Client) UpdateAdditionJobStatus(id string, status jobs.Status) error {
	requestBody := request.Status{
		Status: string(status),
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	url := c.config.jobsAPI
	url.Path = path.Join(additionBasePath, id)

	req, err := http.NewRequest(http.MethodPatch, url.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Add(contentTypeHeaderName, jsonContentType)

	resp, err := c.config.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("could not update job: got status code %d", resp.StatusCode)
	}

	return nil
}
