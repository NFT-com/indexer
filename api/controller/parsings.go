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

type ParsingRepository interface {
	Insert(parsing *jobs.Parsing) (string, error)
	Update(parsing *jobs.Parsing) error
	Retrieve(parsingID string) (*jobs.Parsing, error)
	Find(wheres ...string) ([]*jobs.Parsing, error)
}

type Parsings struct {
	parsings ParsingRepository
	validate *validator.Validate
}

func NewParsings(parsings ParsingRepository) *Parsings {

	p := Parsings{
		parsings: parsings,
		validate: validator.New(),
	}

	return &p
}

func (p *Parsings) Create(ctx echo.Context) error {

	var req []*api.CreateParsingJob
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = p.validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var parsings []*jobs.Parsing
	for _, r := range req {
		parsing := jobs.Parsing{
			ID:          uuid.New().String(),
			ChainID:     r.ChainID,
			Addresses:   r.Addresses,
			EventTypes:  r.EventTypes,
			StartHeight: r.StartHeight,
			EndHeight:   r.EndHeight,
			Status:      jobs.StatusCreated,
			Data:        r.Data,
		}
		parsingID, err := p.parsings.Insert(&parsing)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		parsing.ID = parsingID
		parsings = append(parsings, &parsing)
	}

	return ctx.JSON(http.StatusCreated, parsings)
}

func (p *Parsings) Read(ctx echo.Context) error {

	parsingID := ctx.Param("parsing_id")
	if len(parsingID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing parsing ID")
	}

	parsing, err := p.parsings.Retrieve(parsingID)
	if errors.Is(err, sql.ErrNoRows) {
		return ctx.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, parsing)
}

func (p *Parsings) Update(ctx echo.Context) error {

	parsingID := ctx.Param("parsing_id")
	if len(parsingID) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "missing parsing ID")
	}

	var req api.UpdateParsingJob
	err := ctx.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = p.validate.Struct(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	parsing, err := p.parsings.Retrieve(parsingID)
	if errors.Is(err, sql.ErrNoRows) {
		return ctx.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if parsing.Status == jobs.StatusCreated && req.Status != jobs.StatusQueued ||
		parsing.Status == jobs.StatusQueued && req.Status != jobs.StatusProcessing ||
		parsing.Status == jobs.StatusProcessing && req.Status != jobs.StatusFinished && req.Status != jobs.StatusFailed {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid status transition (%s -> %s)", parsing.Status, req.Status))
	}

	parsing.Status = req.Status
	err = p.parsings.Update(parsing)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, parsing)
}

func (p *Parsings) Index(ctx echo.Context) error {

	var wheres []string

	status := ctx.QueryParam("status")
	if status != "" && !jobs.StatusValid(status) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid job status (%s)", status)
	}
	if status != "" {
		wheres = append(wheres, filters.Eq("status", status))
	}

	parsings, err := p.parsings.Find(wheres...)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, parsings)
}
