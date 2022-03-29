package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
	"github.com/NFT-com/indexer/service/request"
)

// DiscoveryWebsocketConnection handles a new websocket connection.
func (h *Handler) DiscoveryWebsocketConnection(ctx echo.Context) error {
	keys := make(map[string]interface{})

	return h.wsHandler.HandleRequestWithKeys(
		ctx.Response(),
		ctx.Request(),
		broadcaster.WithHandler(keys, broadcaster.DiscoveryHandlerValue),
	)
}

// CreateDiscoveryJob handles the api request to create new discovery job.
func (h *Handler) CreateDiscoveryJob(ctx echo.Context) error {
	var req request.Discovery
	err := ctx.Bind(&req)
	if err != nil {
		return badRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return badRequest(err)
	}

	job := jobs.Discovery{
		ChainURL:     req.ChainURL,
		ChainType:    req.ChainType,
		BlockNumber:  req.BlockNumber,
		Addresses:    req.Addresses,
		StandardType: req.StandardType,
	}

	newJob, err := h.jobs.CreateDiscoveryJob(job)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusCreated, *newJob)
}

// ListDiscoveryJobs handles the api request to retrieve all the discovery jobs.
func (h *Handler) ListDiscoveryJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(statusQueryKey)
	status, err := jobs.ParseStatus(rawStatus)
	if err != nil {
		return badRequest(err)
	}

	jobs, err := h.jobs.ListDiscoveryJobs(status)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, jobs)
}

// GetDiscoveryJob handles the api request to retrieve a discovery job.
func (h *Handler) GetDiscoveryJob(ctx echo.Context) error {
	id := ctx.Param(jobIDParamKey)

	job, err := h.jobs.GetDiscoveryJob(id)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, *job)
}

// UpdateDiscoveryJobStatus handles the api request to update a discovery job status.
func (h *Handler) UpdateDiscoveryJobStatus(ctx echo.Context) error {
	id := ctx.Param(jobIDParamKey)

	var req request.Status
	err := ctx.Bind(&req)
	if err != nil {
		return badRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return badRequest(err)
	}

	newState, err := jobs.ParseStatus(req.Status)
	if err != nil {
		return badRequest(err)
	}

	err = h.jobs.UpdateDiscoveryJobStatus(id, newState)
	if err != nil {
		return internalError(err)
	}

	return ctx.NoContent(http.StatusOK)
}
