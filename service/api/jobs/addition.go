package jobs

import (
	"net/http"

	"github.com/NFT-com/indexer/service/api"
	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
	"github.com/NFT-com/indexer/service/request"
)

// AdditionWebsocketConnection handles a new websocket connection.
func (h *Handler) AdditionWebsocketConnection(ctx echo.Context) error {
	keys := make(map[string]interface{})

	params := ctx.QueryParams()
	if params.Has(statusQueryKey) {
		keys = broadcaster.WithStatus(keys, params.Get(statusQueryKey))
	}

	return h.wsHandler.HandleRequestWithKeys(
		ctx.Response(),
		ctx.Request(),
		broadcaster.WithHandler(keys, broadcaster.AdditionHandlerValue),
	)
}

// CreateAdditionJob handles the api request to create new addition job.
func (h *Handler) CreateAdditionJob(ctx echo.Context) error {
	var req request.Addition
	err := ctx.Bind(&req)
	if err != nil {
		return api.BadRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return api.BadRequest(err)
	}

	job := jobs.Addition{
		ChainURL:     req.ChainURL,
		ChainID:      req.ChainID,
		ChainType:    req.ChainType,
		BlockNumber:  req.BlockNumber,
		Address:      req.Address,
		StandardType: req.StandardType,
		TokenID:      req.TokenID,
	}

	newJob, err := h.jobs.CreateAdditionJob(job)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.JSON(http.StatusCreated, *newJob)
}

// CreateAdditionJobs handles the api request to create multiple new addition jobs.
func (h *Handler) CreateAdditionJobs(ctx echo.Context) error {
	var req request.Additions
	err := ctx.Bind(&req)
	if err != nil {
		return api.BadRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return api.BadRequest(err)
	}

	jobList := make([]jobs.Addition, 0, len(req.Jobs))
	for _, j := range req.Jobs {
		job := jobs.Addition{
			ChainURL:     j.ChainURL,
			ChainID:      j.ChainID,
			ChainType:    j.ChainType,
			BlockNumber:  j.BlockNumber,
			Address:      j.Address,
			StandardType: j.StandardType,
			TokenID:      j.TokenID,
		}

		jobList = append(jobList, job)
	}

	err = h.jobs.CreateAdditionJobs(jobList)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.NoContent(http.StatusCreated)
}

// ListAdditionJobs handles the api request to retrieve all the addition jobs.
func (h *Handler) ListAdditionJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(statusQueryKey)
	status, err := jobs.ParseStatus(rawStatus)
	if err != nil {
		return api.BadRequest(err)
	}

	jobs, err := h.jobs.ListAdditionJobs(status)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.JSON(http.StatusOK, jobs)
}

// GetAdditionJob handles the api request to retrieve a discovery job.
func (h *Handler) GetAdditionJob(ctx echo.Context) error {
	id := ctx.Param(jobIDParamKey)

	job, err := h.jobs.GetAdditionJob(id)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.JSON(http.StatusOK, *job)
}

// UpdateAdditionJobStatus handles the api request to update a addition job status.
func (h *Handler) UpdateAdditionJobStatus(ctx echo.Context) error {
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

	err = h.jobs.UpdateAdditionJobStatus(id, newState)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.NoContent(http.StatusOK)
}
