package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
	"github.com/NFT-com/indexer/service/request"
)

// ParsingWebsocketConnection handles a new websocket connection.
func (h *Handler) ParsingWebsocketConnection(ctx echo.Context) error {
	keys := make(map[string]interface{})

	return h.wsHandler.HandleRequestWithKeys(
		ctx.Response(),
		ctx.Request(),
		broadcaster.WithHandler(keys, broadcaster.ParsingHandlerValue),
	)
}

// CreateParsingJob handles the api request to create new parsing job.
func (h *Handler) CreateParsingJob(ctx echo.Context) error {
	var req request.Parsing
	err := ctx.Bind(&req)
	if err != nil {
		return badRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return badRequest(err)
	}

	job := jobs.Parsing{
		ChainURL:     req.ChainURL,
		ChainType:    req.ChainType,
		BlockNumber:  req.BlockNumber,
		Address:      req.Address,
		StandardType: req.StandardType,
		EventType:    req.EventType,
	}

	newJob, err := h.jobs.CreateParsingJob(job)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusCreated, *newJob)
}

// ListParsingJobs handles the api request to retrieve all the parsing jobs.
func (h *Handler) ListParsingJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(statusQueryKey)

	status, err := jobs.ParseStatus(rawStatus)
	if err != nil {
		return badRequest(err)
	}

	jobs, err := h.jobs.ListParsingJobs(status)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, jobs)
}

// GetParsingJob handles the api request to retrieve a parsing job.
func (h *Handler) GetParsingJob(ctx echo.Context) error {
	id := ctx.Param(jobIDParamKey)

	job, err := h.jobs.GetParsingJob(id)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, *job)
}

// UpdateParsingJobStatus handles the api request to update a parsing job status.
func (h *Handler) UpdateParsingJobStatus(ctx echo.Context) error {
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

	err = h.jobs.UpdateParsingJobStatus(id, newState)
	if err != nil {
		return internalError(err)
	}

	return ctx.NoContent(http.StatusOK)
}
