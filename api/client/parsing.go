package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/NFT-com/indexer/models/api"
	"github.com/NFT-com/indexer/models/jobs"
)

func (c *Client) CreateParsingJobs(parsings []*jobs.Parsing) error {

	reqs := make([]*api.CreateParsingJob, 0, len(parsings))
	for _, parsing := range parsings {
		req := api.CreateParsingJob{
			ChainID:     parsing.ChainID,
			Addresses:   parsing.Addresses,
			EventTypes:  parsing.EventTypes,
			StartHeight: parsing.StartHeight,
			EndHeight:   parsing.EndHeight,
			Data:        parsing.Data,
		}

		reqs = append(reqs, &req)
	}
	body, err := json.Marshal(reqs)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	url := *c.url
	url.Path = path.Join(parsingBasePath, "batch")

	_, err = c.client.Post(url.String(), jsonContentType, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}

	return nil
}

func (c *Client) ListParsingJobs(status string) ([]jobs.Parsing, error) {

	params := url.Values{}
	params.Set("status", string(status))

	url := *c.url
	url.Path = parsingBasePath
	url.RawQuery = params.Encode()

	res, err := c.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("could not list job: got status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}
	defer res.Body.Close()

	parsings := make([]jobs.Parsing, 0)
	err = json.Unmarshal(body, &parsings)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return parsings, nil
}

func (c *Client) GetParsingJob(parsingID string) (*jobs.Parsing, error) {

	url := *c.url
	url.Path = path.Join(parsingBasePath, parsingID)

	res, err := c.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("could not get job: got status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}
	defer res.Body.Close()

	var job jobs.Parsing
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

	url := *c.url
	url.Path = path.Join(parsingBasePath, "highest")
	url.RawQuery = params.Encode()

	res, err := c.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("could not highest block number job: got status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}
	defer res.Body.Close()

	var job jobs.Parsing
	err = json.Unmarshal(body, &job)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &job, nil
}

func (c *Client) UpdateParsingJobStatus(id string, status string) error {

	req := api.UpdateParsingJob{
		Status: string(status),
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	url := *c.url
	url.Path = path.Join(parsingBasePath, id)

	rawReq, err := http.NewRequest(http.MethodPatch, url.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	rawReq.Header.Add(contentTypeHeaderName, jsonContentType)

	res, err := c.client.Do(rawReq)
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("could not update job: got status code %d", res.StatusCode)
	}

	return nil
}
