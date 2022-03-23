package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/request"
)

func (c *Client) SubscribeNewDiscoveryJob(discoveryJobs chan jobs.Discovery) error {
	requestURL := fmt.Sprintf("%s/ws/%s", c.options.websocketURL.String(), discoveryBasePath)
	connection, _, err := c.wsClient.Dial(requestURL, nil)
	if err != nil {
		return fmt.Errorf("could not dial to websocket: %w", err)
	}

	go func() {
		for {
			select {
			case <-c.close:
				return
			default:
			}

			newDiscoveryJob := jobs.Discovery{}
			err := connection.ReadJSON(&newDiscoveryJob)
			if err != nil {
				c.log.Error().Err(err).Msg("could not read message socket")
				continue
			}

			discoveryJobs <- newDiscoveryJob
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

	requestURL := fmt.Sprintf("%s/%s", c.options.httpURL.String(), discoveryBasePath)
	resp, err := c.httpClient.Post(requestURL, jsonContentType, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
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
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s?status=%s", c.options.httpURL.String(), discoveryBasePath, status))
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
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
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), discoveryBasePath, id))
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
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

	requestURL := fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), discoveryBasePath, id)
	req, err := http.NewRequest(http.MethodPatch, requestURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Add(contentTypeHeaderName, jsonContentType)

	_, err = c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}

	return nil
}

func (c *Client) RequeueDiscoveryJob(id string) (*jobs.Discovery, error) {
	resp, err := c.httpClient.Post(fmt.Sprintf("%s/%s/%s/requeue", c.options.httpURL.String(), discoveryBasePath, id), jsonContentType, nil)
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
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
