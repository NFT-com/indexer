package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/models/api"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/storage/filters"
)

type ActionRepository interface {
	Insert(action *jobs.Action) (string, error)
	Update(action *jobs.Action) error
	Retrieve(actionID string) (*jobs.Action, error)
	Find(wheres ...string) ([]*jobs.Action, error)
}

type Actions struct {
	actions  ActionRepository
	validate *validator.Validate
}

func NewActions(actions ActionRepository) *Actions {

	a := Actions{
		actions:  actions,
		validate: validator.New(),
	}

	return &a
}

func (a *Actions) Create(ctx echo.Context) error {

	var req []*api.CreateActionJob
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = a.validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var actions []*jobs.Action
	for _, r := range req {
		action := jobs.Action{
			ID:         uuid.New().String(),
			ChainID:    r.ChainID,
			Address:    r.Address,
			TokenID:    r.TokenID,
			ActionType: r.ActionType,
			Height:     r.Height,
			Data:       r.Data,
			Status:     jobs.StatusCreated,
		}
		actionID, err := a.actions.Insert(&action)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		action.ID = actionID
		actions = append(actions, &action)
	}

	return ctx.JSON(http.StatusCreated, actions)
}

func (a *Actions) Read(ctx echo.Context) error {

	actionID := ctx.Param("action_id")
	if len(actionID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing action ID")
	}

	action, err := a.actions.Retrieve(actionID)
	if errors.Is(err, sql.ErrNoRows) {
		return ctx.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, action)
}

func (a *Actions) Update(ctx echo.Context) error {

	actionID := ctx.Param("action_id")
	if len(actionID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing action ID")
	}

	var req api.UpdateActionJob
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = a.validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	action, err := a.actions.Retrieve(actionID)
	if errors.Is(err, sql.ErrNoRows) {
		return ctx.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if action.Status == jobs.StatusCreated && req.Status != jobs.StatusQueued ||
		action.Status == jobs.StatusQueued && req.Status != jobs.StatusProcessing ||
		action.Status == jobs.StatusProcessing && req.Status != jobs.StatusFinished && req.Status != jobs.StatusFailed {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid status transition (%s -> %s)", action.Status, req.Status))
	}

	action.Status = req.Status
	err = a.actions.Update(action)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, action)
}

func (a *Actions) Index(ctx echo.Context) error {

	var wheres []string

	status := ctx.QueryParam("status")
	if status != "" && !jobs.StatusValid(status) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid job status (%s)", status)
	}
	if status != "" {
		wheres = append(wheres, filters.Eq("status", status))
	}

	actions, err := a.actions.Find(wheres...)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, actions)
}
