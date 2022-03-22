package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
	"github.com/NFT-com/indexer/service/request"
)

func (h *Handler) NewParsingWebsocketConnection(ctx echo.Context) error {
	return h.wsHandler.HandleRequestWithKeys(
		ctx.Response(),
		ctx.Request(),
		broadcaster.Keys{}.
			WithHandler(broadcaster.ParsingHandlerValue),
	)
}

func (h *Handler) CreateParsingJob(ctx echo.Context) error {
	var req request.Parsing
	err := ctx.Bind(&req)
	if err != nil {
		return unpackError(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return validateError(err)
	}

	job := jobs.Parsing{
		ChainURL:     req.ChainURL,
		ChainType:    req.ChainType,
		BlockNumber:  req.BlockNumber,
		Address:      req.Address,
		StandardType: req.InterfaceType,
		EventType:    req.EventType,
	}

	newJob, err := h.jobs.CreateParsingJob(job)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusCreated, *newJob)
}

func (h *Handler) ListParsingJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(statusQueryKey)

	status, err := jobs.ParseStatus(rawStatus)
	if err != nil {
		return parsingError(err)
	}

	jobs, err := h.jobs.ListParsingJobs(status)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, jobs)
}

func (h *Handler) GetParsingJob(ctx echo.Context) error {
	jobID := ctx.Param(jobIDParamKey)

	job, err := h.jobs.GetParsingJob(jobID)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, *job)
}

func (h *Handler) UpdateParsingJobStatus(ctx echo.Context) error {
	jobID := ctx.Param(jobIDParamKey)

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

	err = h.jobs.UpdateParsingJobState(jobID, newState)
	if err != nil {
		return apiError(err)
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *Handler) RequeueParsingJob(ctx echo.Context) error {
	jobID := ctx.Param(jobIDParamKey)

	job, err := h.jobs.RequeueParsingJob(jobID)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusCreated, *job)
}
