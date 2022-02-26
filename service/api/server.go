package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/queue"
	"github.com/NFT-com/indexer/queue/producer"
	"github.com/NFT-com/indexer/service/request"
)

type API struct {
	discoveryJobQueueName string
	parseJobQueueName     string
	produce               *producer.Producer
}

func NewAPI(discoveryJobQueueName string, parseJobQueueName string, produce *producer.Producer) (*API, error) {
	a := API{
		discoveryJobQueueName: discoveryJobQueueName,
		parseJobQueueName:     parseJobQueueName,
		produce:               produce,
	}

	return &a, nil
}

func (a *API) PublishDiscoveryJob(ctx echo.Context) error {
	var req request.DiscoveryJob
	if err := ctx.Bind(&req); err != nil {
		return unpackError(err)
	}

	job := queue.DiscoveryJob{
		ChainURL:   req.ChainURL,
		ChainType:  req.ChainType,
		StartIndex: req.StartIndex,
		EndIndex:   req.EndIndex,
		Contracts:  req.Contracts,
	}
	fmt.Println(job)
	if err := a.produce.PublishDiscoveryJob(a.discoveryJobQueueName, job); err != nil {
		return apiError(err)
	}

	return ctx.String(http.StatusOK, http.StatusText(http.StatusOK))
}

func (a *API) PublishParseJob(ctx echo.Context) error {
	var req request.ParseJob
	if err := ctx.Bind(&req); err != nil {
		return unpackError(err)
	}

	job := queue.ParseJob{
		ID:              req.ID,
		NetworkID:       req.NetworkID,
		ChainID:         req.ChainID,
		Block:           req.Block,
		TransactionHash: req.TransactionHash,
		Address:         req.Address,
		AddressType:     req.AddressType,
		Topic:           req.Topic,
		IndexedData:     req.IndexedData,
		Data:            req.Data,
	}
	if err := a.produce.PublishParseJob(a.parseJobQueueName, job); err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, nil)
}
