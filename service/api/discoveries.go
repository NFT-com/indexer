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
	rawStatus := ctx.QueryParam(StatusQueryKey)
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
	rawJobID := ctx.Param(jobIDParamKey)
	jobID := jobs.ID(rawJobID)

	job, err := h.jobs.GetDiscoveryJob(jobID)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, *job)
}

func (h *Handler) UpdateDiscoveryJobStatus(ctx echo.Context) error {
	rawJobID := ctx.Param(jobIDParamKey)
	jobID := jobs.ID(rawJobID)

	var req request.Status
	err := ctx.Bind(&req)
	if err != nil {
		return unpackError(err)
	}

	newState, err := jobs.ParseStatus(req.Status)
	if err != nil {
		return parsingError(err)
	}

	err = h.jobs.UpdateDiscoveryJobState(jobID, newState)
	if err != nil {
		return apiError(err)
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) RequeueDiscoveryJob(ctx echo.Context) error {
	rawJobID := ctx.Param(jobIDParamKey)
	jobID := jobs.ID(rawJobID)

	newJob, err := h.jobs.RequeueDiscoveryJob(jobID)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusCreated, *newJob)
}
