package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
	"github.com/NFT-com/indexer/service/request"
)

func (h *Handler) NewDiscoveryWebsocketConnection(ctx echo.Context) error {
	return h.wsHandler.HandleRequestWithKeys(
		ctx.Response(),
		ctx.Request(),
		broadcaster.Keys{}.
			WithHandler(broadcaster.DiscoveryHandlerValue),
	)
}

func (h *Handler) CreateDiscoveryJob(ctx echo.Context) error {
	var req request.Discovery
	err := ctx.Bind(&req)
	if err != nil {
		return unpackError(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return validateError(err)
	}

	job := jobs.Discovery{
		ChainURL:     req.ChainURL,
		ChainType:    req.ChainType,
		BlockNumber:  req.BlockNumber,
		Addresses:    req.Addresses,
		StandardType: req.InterfaceType,
	}

	newJob, err := h.jobs.CreateDiscoveryJob(job)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusCreated, *newJob)
}

func (h *Handler) ListDiscoveryJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(statusQueryKey)
	status, err := jobs.ParseStatus(rawStatus)
	if err != nil {
		return parsingError(err)
	}

	jobs, err := h.jobs.ListDiscoveryJobs(status)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, jobs)
}

func (h *Handler) GetDiscoveryJob(ctx echo.Context) error {
	id := ctx.Param(jobIDParamKey)

	job, err := h.jobs.GetDiscoveryJob(id)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, *job)
}

func (h *Handler) UpdateDiscoveryJobStatus(ctx echo.Context) error {
	id := ctx.Param(jobIDParamKey)

	var req request.Status
	err := ctx.Bind(&req)
	if err != nil {
		return unpackError(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return validateError(err)
	}

	newState, err := jobs.ParseStatus(req.Status)
	if err != nil {
		return parsingError(err)
	}

	err = h.jobs.UpdateDiscoveryJobState(id, newState)
	if err != nil {
		return apiError(err)
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) RequeueDiscoveryJob(ctx echo.Context) error {
	id := ctx.Param(jobIDParamKey)

	newJob, err := h.jobs.RequeueDiscoveryJob(id)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusCreated, *newJob)
}
