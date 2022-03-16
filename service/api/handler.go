package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/olahol/melody.v1"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/service/broadcaster"
	"github.com/NFT-com/indexer/service/request"
)

const (
	DiscoveryJobIDParamKey = "discovery_id"
	ParsingJobIDParamKey   = "parsing_id"

	StatusQueryKey = "status"
)

type Handler struct {
	wsHandler     *melody.Melody
	jobController JobController
}

func NewHandler(wsHandler *melody.Melody, jobController JobController) *Handler {
	s := Handler{
		wsHandler:     wsHandler,
		jobController: jobController,
	}

	return &s
}

func (h *Handler) ApplyRoutes(server *echo.Echo) {
	websocketGroup := server.Group("/ws")
	{
		websocketGroup.GET("/discoveries", h.NewDiscoveryWebsocketConnection)
		websocketGroup.GET("/parsings", h.NewParsingWebsocketConnection)
	}

	discoveriesJobGroup := server.Group("/discoveries")
	{
		discoveriesJobGroup.POST("", h.CreateDiscoveryJob)
		discoveriesJobGroup.GET("", h.ListDiscoveryJobs)
		discoveriesJobGroup.GET("/:"+DiscoveryJobIDParamKey, h.GetDiscoveryJob)
		discoveriesJobGroup.PATCH("/:"+DiscoveryJobIDParamKey, h.UpdateDiscoveryJobStatus)
		discoveriesJobGroup.POST("/:"+DiscoveryJobIDParamKey+"/requeue", h.RequeueDiscoveryJob)
	}

	parsingsJobGroup := server.Group("/parsings")
	{
		parsingsJobGroup.POST("", h.CreateParsingJob)
		parsingsJobGroup.GET("", h.ListParsingJobs)
		parsingsJobGroup.GET("/:"+ParsingJobIDParamKey, h.GetParsingJob)
		parsingsJobGroup.PATCH("/:"+ParsingJobIDParamKey, h.UpdateParsingJobStatus)
		parsingsJobGroup.POST("/:"+ParsingJobIDParamKey+"/requeue", h.RequeueParsingJob)
	}
}

func (h *Handler) NewDiscoveryWebsocketConnection(ctx echo.Context) error {
	return h.wsHandler.HandleRequestWithKeys(
		ctx.Response(),
		ctx.Request(),
		broadcaster.Keys{}.
			WithHandler(broadcaster.DiscoveryHandlerValue),
	)
}

func (h *Handler) NewParsingWebsocketConnection(ctx echo.Context) error {
	return h.wsHandler.HandleRequestWithKeys(
		ctx.Response(),
		ctx.Request(),
		broadcaster.Keys{}.
			WithHandler(broadcaster.ParsingHandlerValue),
	)
}

func (h *Handler) CreateDiscoveryJob(ctx echo.Context) error {
	var req request.Discovery
	if err := ctx.Bind(&req); err != nil {
		return unpackError(err)
	}

	job := jobs.Discovery{
		ChainURL:     req.ChainURL,
		ChainType:    req.ChainType,
		BlockNumber:  req.BlockNumber,
		Addresses:    req.Addresses,
		StandardType: req.InterfaceType,
	}

	newJob, err := h.jobController.CreateDiscoveryJob(job)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, *newJob)
}

func (h *Handler) ListDiscoveryJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(StatusQueryKey)
	status, err := jobs.ParseStatus(rawStatus)
	if err != nil {
		return parsingError(err)
	}

	jobs, err := h.jobController.ListDiscoveryJobs(status)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, jobs)
}

func (h *Handler) GetDiscoveryJob(ctx echo.Context) error {
	jobID := ctx.Param(DiscoveryJobIDParamKey)

	discoveryJob, err := h.jobController.GetDiscoveryJob(jobID)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, *discoveryJob)
}

func (h *Handler) UpdateDiscoveryJobStatus(ctx echo.Context) error {
	jobID := ctx.Param(DiscoveryJobIDParamKey)

	var req request.Status
	if err := ctx.Bind(&req); err != nil {
		return unpackError(err)
	}

	newState, err := jobs.ParseStatus(req.Status)
	if err != nil {
		return parsingError(err)
	}

	err = h.jobController.UpdateDiscoveryJobState(jobID, newState)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) RequeueDiscoveryJob(ctx echo.Context) error {
	jobID := ctx.Param(DiscoveryJobIDParamKey)

	newJob, err := h.jobController.RequeueDiscoveryJob(jobID)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, *newJob)
}

func (h *Handler) CreateParsingJob(ctx echo.Context) error {
	var req request.Parsing
	if err := ctx.Bind(&req); err != nil {
		return unpackError(err)
	}

	job := jobs.Parsing{
		ChainURL:     req.ChainURL,
		ChainType:    req.ChainType,
		BlockNumber:  req.BlockNumber,
		Address:      req.Address,
		StandardType: req.InterfaceType,
		EventType:    req.EventType,
	}

	newJob, err := h.jobController.CreateParsingJob(job)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, *newJob)
}

func (h *Handler) ListParsingJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(StatusQueryKey)
	status, err := jobs.ParseStatus(rawStatus)
	if err != nil {
		return parsingError(err)
	}

	jobs, err := h.jobController.ListParsingJobs(status)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, jobs)
}

func (h *Handler) GetParsingJob(ctx echo.Context) error {
	jobID := ctx.Param(ParsingJobIDParamKey)

	job, err := h.jobController.GetParsingJob(jobID)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, *job)
}

func (h *Handler) UpdateParsingJobStatus(ctx echo.Context) error {
	jobID := ctx.Param(ParsingJobIDParamKey)

	var req request.Status
	if err := ctx.Bind(&req); err != nil {
		return unpackError(err)
	}

	newState, err := jobs.ParseStatus(req.Status)
	if err != nil {
		return parsingError(err)
	}

	err = h.jobController.UpdateParsingJobState(jobID, newState)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) RequeueParsingJob(ctx echo.Context) error {
	jobID := ctx.Param(ParsingJobIDParamKey)

	job, err := h.jobController.RequeueParsingJob(jobID)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, *job)
}
