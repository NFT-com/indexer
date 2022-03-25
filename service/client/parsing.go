package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/request"
)

func (c *Client) SubscribeNewParsingJob(parsingJobs chan jobs.Parsing) error {
	requestURL := fmt.Sprintf("%s/ws/%s", c.options.websocketURL.String(), parsingBasePath)
	connection, _, err := c.options.wsDialer.Dial(requestURL, nil)
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

			job := jobs.Parsing{}
			err := connection.ReadJSON(&job)
			if err != nil {
				c.log.Error().Err(err).Msg("could not read message socket")
				continue
			}

			parsingJobs <- job
		}
	}()

	return nil
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
		return nil, fmt.Errorf("could not marshal request: %w", err)
	}

	requestURL := fmt.Sprintf("%s/%s", c.options.httpURL.String(), parsingBasePath)
	resp, err := c.options.httpClient.Post(requestURL, jsonContentType, bytes.NewReader(body))
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

	newJob := jobs.Parsing{}
	err = json.Unmarshal(responseBody, &newJob)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &newJob, nil
}

func (c *Client) ListParsingJobs(status jobs.Status) ([]jobs.Parsing, error) {
	resp, err := c.options.httpClient.Get(fmt.Sprintf("%s/%s?status=%s", c.options.httpURL.String(), parsingBasePath, status))
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

	jobList := make([]jobs.Parsing, 0)
	err = json.Unmarshal(body, &jobList)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return jobList, nil
}

func (c *Client) GetParsingJob(id string) (*jobs.Parsing, error) {
	resp, err := c.options.httpClient.Get(fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), parsingBasePath, id))
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

	job := jobs.Parsing{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &job, nil
}

func (c *Client) UpdateParsingJobState(id string, status jobs.Status) error {
	requestBody := request.Status{
		Status: string(status),
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	requestURL := fmt.Sprintf("%s/%s/%s", c.options.httpURL.String(), parsingBasePath, id)
	req, err := http.NewRequest(http.MethodPatch, requestURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Add(contentTypeHeaderName, jsonContentType)

	_, err = c.options.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}

	return nil
}

func (c *Client) RequeueParsingJob(id string) (*jobs.Parsing, error) {
	resp, err := c.options.httpClient.Post(fmt.Sprintf("%s/%s/%s/requeue", c.options.httpURL.String(), parsingBasePath, id), jsonContentType, nil)
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

	job := jobs.Parsing{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &job, nil
}
