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

func (c *Client) CreateActionJob(action *jobs.Action) (*jobs.Action, error) {

	req := api.CreateActionJob{
		ChainURL:    action.ChainURL,
		ChainID:     action.ChainID,
		ChainType:   action.ChainType,
		BlockNumber: action.BlockNumber,
		Address:     action.Address,
		Standard:    action.Standard,
		TokenID:     action.TokenID,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request: %w", err)
	}

	url := *c.url
	url.Path = actionBasePath

	res, err := c.client.Post(url.String(), jsonContentType, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not perform request: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("could not create job: got status code %d", res.StatusCode)
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body: %w", err)
	}
	defer res.Body.Close()

	var newJob jobs.Action
	err = json.Unmarshal(responseBody, &newJob)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &newJob, nil
}

func (c *Client) CreateActionJobs(actions []*jobs.Action) error {

	reqs := make([]api.CreateActionJob, 0, len(actions))
	for _, action := range actions {
		req := api.CreateActionJob{
			ChainURL:    action.ChainURL,
			ChainID:     action.ChainID,
			ChainType:   action.ChainType,
			BlockNumber: action.BlockNumber,
			Address:     action.Address,
			Standard:    action.Standard,
			TokenID:     action.TokenID,
		}
		reqs = append(reqs, req)
	}

	body, err := json.Marshal(reqs)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	url := *c.url
	url.Path = path.Join(actionBasePath, "batch")

	_, err = c.client.Post(url.String(), jsonContentType, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not perform request: %w", err)
	}

	return nil
}

func (c *Client) ListActionJobs(status string) ([]*jobs.Action, error) {

	params := url.Values{}
	params.Set("status", string(status))

	url := *c.url
	url.Path = actionBasePath
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

	var actions []*jobs.Action
	err = json.Unmarshal(body, &actions)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return actions, nil
}

func (c *Client) GetActionJob(actionID string) (*jobs.Action, error) {

	url := *c.url
	url.Path = path.Join(actionBasePath, actionID)

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

	var action jobs.Action
	err = json.Unmarshal(body, &action)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &action, nil
}

func (c *Client) UpdateActionJobStatus(actionID string, status string) error {

	req := api.UpdateActionJob{
		Status: status,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	url := *c.url
	url.Path = path.Join(actionBasePath, actionID)

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
