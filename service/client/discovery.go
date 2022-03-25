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

func (c *Client) SubscribeNewDiscoveryJob(discoveryJobs chan jobs.Discovery) error {
	url := c.options.jobsWebsocket
	url.Path = path.Join("ws", discoveryBasePath)

	connection, _, err := c.options.dialer.Dial(url.String(), nil)
	if err != nil {
		return fmt.Errorf("could not dial websocket: %w", err)
	}

	go func() {
		for {
			select {
			case <-c.close:
				return
			default:
			}

			job := jobs.Discovery{}
			err := connection.ReadJSON(&job)
			if err != nil {
				c.log.Error().Err(err).Msg("could not read message socket")
				continue
			}

			discoveryJobs <- job
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
		return nil, fmt.Errorf("could not marshal request: %w", err)
	}

	url := c.options.jobsAPI
	url.Path = discoveryBasePath

	resp, err := c.options.client.Post(url.String(), jsonContentType, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close response body: %w", err)
	}

	newJob := jobs.Discovery{}
	err = json.Unmarshal(responseBody, &newJob)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &newJob, nil
}

func (c *Client) ListDiscoveryJobs(status jobs.Status) ([]jobs.Discovery, error) {
	params := url.Values{}
	params.Set("status", string(status))

	url := c.options.jobsAPI
	url.Path = discoveryBasePath
	url.RawQuery = params.Encode()

	resp, err := c.options.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close response body: %w", err)
	}

	jobList := make([]jobs.Discovery, 0)
	err = json.Unmarshal(body, &jobList)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return jobList, nil
}

func (c *Client) GetDiscoveryJob(id string) (*jobs.Discovery, error) {
	url := c.options.jobsAPI
	url.Path = path.Join(discoveryBasePath, id)

	resp, err := c.options.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not close response body: %w", err)
	}

	job := jobs.Discovery{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &job, nil
}

func (c *Client) UpdateDiscoveryJobState(id string, status jobs.Status) error {
	requestBody := request.Status{
		Status: string(status),
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	url := c.options.jobsAPI
	url.Path = path.Join(discoveryBasePath, id)

	req, err := http.NewRequest(http.MethodPatch, url.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Add(contentTypeHeaderName, jsonContentType)

	_, err = c.options.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}

	return nil
}
