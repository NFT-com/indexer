package jobs

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/api"
	"github.com/NFT-com/indexer/service/broadcaster"
	"github.com/NFT-com/indexer/service/request"
)

// ParsingWebsocketConnection handles a new websocket connection.
func (h *Handler) ParsingWebsocketConnection(ctx echo.Context) error {
	keys := make(map[string]interface{})

	params := ctx.QueryParams()
	if params.Has(statusQueryKey) {
		keys = broadcaster.WithStatus(keys, params.Get(statusQueryKey))
	}

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
		return api.BadRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return api.BadRequest(err)
	}

	job := jobs.Parsing{
		ChainURL:     req.ChainURL,
		ChainID:      req.ChainID,
		ChainType:    req.ChainType,
		BlockNumber:  req.BlockNumber,
		Address:      req.Address,
		StandardType: req.StandardType,
		EventType:    req.EventType,
	}

	newJob, err := h.jobs.CreateParsingJob(job)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.JSON(http.StatusCreated, *newJob)
}

// CreateParsingJobs handles the api request to create multiple new parsing jobs.
func (h *Handler) CreateParsingJobs(ctx echo.Context) error {
	var req request.Parsings
	err := ctx.Bind(&req)
	if err != nil {
		return api.BadRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return api.BadRequest(err)
	}

	jobList := make([]jobs.Parsing, 0, len(req.Jobs))
	for _, j := range req.Jobs {
		job := jobs.Parsing{
			ChainURL:     j.ChainURL,
			ChainID:      j.ChainID,
			ChainType:    j.ChainType,
			BlockNumber:  j.BlockNumber,
			Address:      j.Address,
			StandardType: j.StandardType,
			EventType:    j.EventType,
		}

		jobList = append(jobList, job)
	}

	err = h.jobs.CreateParsingJobs(jobList)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.NoContent(http.StatusCreated)
}

// ListParsingJobs handles the api request to retrieve all the parsing jobs.
func (h *Handler) ListParsingJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(statusQueryKey)

	status, err := jobs.ParseStatus(rawStatus)
	if err != nil {
		return api.BadRequest(err)
	}

	jobs, err := h.jobs.ListParsingJobs(status)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.JSON(http.StatusOK, jobs)
}

// GetParsingJob handles the api request to retrieve a parsing job.
func (h *Handler) GetParsingJob(ctx echo.Context) error {
	id := ctx.Param(jobIDParamKey)

	job, err := h.jobs.GetParsingJob(id)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.JSON(http.StatusOK, *job)
}

// GetHighestBlockNumberParsingJob handles the api request to retrieve the highest parsing block number.
func (h *Handler) GetHighestBlockNumberParsingJob(ctx echo.Context) error {
	var (
		chainURL     = ctx.QueryParam(chainURLQueryKey)
		chainType    = ctx.QueryParam(chainTypeQueryKey)
		address      = ctx.QueryParam(addressQueryKey)
		standardType = ctx.QueryParam(standardTypeQueryKey)
		eventType    = ctx.QueryParam(eventTypeQueryKey)
	)

	job, err := h.jobs.GetHighestBlockNumberParsingJob(chainURL, chainType, address, standardType, eventType)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.JSON(http.StatusOK, *job)
}

// UpdateParsingJobStatus handles the api request to update a parsing job status.
func (h *Handler) UpdateParsingJobStatus(ctx echo.Context) error {
	id := ctx.Param(jobIDParamKey)

	var req request.Status
	err := ctx.Bind(&req)
	if err != nil {
		return api.BadRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return api.BadRequest(err)
	}

	newState, err := jobs.ParseStatus(req.Status)
	if err != nil {
		return api.BadRequest(err)
	}

	err = h.jobs.UpdateParsingJobStatus(id, newState)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.NoContent(http.StatusOK)
}
