package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/request"
)

// CreateActionJob handles the api request to create new action job.
func (h *Handler) CreateActionJob(ctx echo.Context) error {
	var req request.Action
	err := ctx.Bind(&req)
	if err != nil {
		return badRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return badRequest(err)
	}

	job := jobs.Action{
		ChainURL:    req.ChainURL,
		ChainID:     req.ChainID,
		ChainType:   req.ChainType,
		BlockNumber: req.BlockNumber,
		Address:     req.Address,
		Standard:    req.Standard,
		Event:       req.Event,
		TokenID:     req.TokenID,
		ToAddress:   req.ToAddress,
		Type:        req.Type,
	}

	newJob, err := h.jobs.CreateActionJob(&job)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusCreated, *newJob)
}

// CreateActionJobs handles the api request to create multiple new action jobs.
func (h *Handler) CreateActionJobs(ctx echo.Context) error {
	var req request.Actions
	err := ctx.Bind(&req)
	if err != nil {
		return badRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return badRequest(err)
	}

	jobList := make([]*jobs.Action, 0, len(req.Jobs))
	for _, j := range req.Jobs {
		job := jobs.Action{
			ChainURL:    j.ChainURL,
			ChainID:     j.ChainID,
			ChainType:   j.ChainType,
			BlockNumber: j.BlockNumber,
			Address:     j.Address,
			Standard:    j.Standard,
			Event:       j.Event,
			TokenID:     j.TokenID,
			ToAddress:   j.ToAddress,
			Type:        j.Type,
		}

		jobList = append(jobList, &job)
	}

	err = h.jobs.CreateActionJobs(jobList)
	if err != nil {
		return internalError(err)
	}

	return ctx.NoContent(http.StatusCreated)
}

// ListActionJobs handles the api request to retrieve all the action jobs.
func (h *Handler) ListActionJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(statusQueryKey)
	status, err := jobs.ParseStatus(rawStatus)
	if err != nil {
		return badRequest(err)
	}

	jobs, err := h.jobs.ListActionJobs(status)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, jobs)
}

// GetActionJob handles the api request to retrieve an action job.
func (h *Handler) GetActionJob(ctx echo.Context) error {
	id := ctx.Param(jobIDParamKey)

	job, err := h.jobs.GetActionJob(id)
	if err != nil {
		return internalError(err)
	}

	return ctx.JSON(http.StatusOK, *job)
}

// UpdateActionJobStatus handles the api request to update an action job status.
func (h *Handler) UpdateActionJobStatus(ctx echo.Context) error {
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

	err = h.jobs.UpdateActionJobStatus(id, newState)
	if err != nil {
		return internalError(err)
	}

	return ctx.NoContent(http.StatusOK)
}
