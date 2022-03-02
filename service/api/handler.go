package api

import (
	"github.com/NFT-com/indexer/job"
	"github.com/NFT-com/indexer/service/request"
	"github.com/google/uuid"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	DiscoveryJobIDParamKey = "discovery_id"
	ParsingJobIDParamKey   = "parsing_id"

	StatusQueryKey = "status"
)

type Handler struct {
	discoveryJobsStore DiscoveryJobsStore
	parsingJobsStore   ParsingJobsStore
}

func NewHandler(discoveryJobsStore DiscoveryJobsStore, parsingJobsStore ParsingJobsStore) (*Handler, error) {
	s := Handler{
		discoveryJobsStore: discoveryJobsStore,
		parsingJobsStore:   parsingJobsStore,
	}

	return &s, nil
}

func (h *Handler) ApplyRoutes(server *echo.Echo) {
	deliveryJobGroup := server.Group("/deliveries")
	{
		deliveryJobGroup.POST("", h.CreateDiscoveryJob)
		deliveryJobGroup.GET("", h.ListDiscoveryJobs)
		deliveryJobGroup.GET("/"+DiscoveryJobIDParamKey, h.GetDiscoveryJob)
		deliveryJobGroup.DELETE("/"+DiscoveryJobIDParamKey, h.GetDiscoveryJob)
		deliveryJobGroup.POST("/"+DiscoveryJobIDParamKey+"/requeue", h.RequeueDiscoveryJob)
	}

	parsingsJobGroup := server.Group("/parsings")
	{
		parsingsJobGroup.POST("", h.CreateParsingJob)
		parsingsJobGroup.GET("", h.ListParsingJobs)
		parsingsJobGroup.GET("/"+ParsingJobIDParamKey, h.GetParsingJob)
		parsingsJobGroup.DELETE("/"+ParsingJobIDParamKey, h.CancelParsingJob)
		parsingsJobGroup.POST("/"+ParsingJobIDParamKey+"/requeue", h.Noop)
	}
}

func (h *Handler) CreateDiscoveryJob(ctx echo.Context) error {
	var req request.Discovery
	if err := ctx.Bind(&req); err != nil {
		return unpackError(err)
	}

	discoveryJob := job.Discovery{
		ID:            uuid.New().String(),
		ChainURL:      req.ChainURL,
		ChainType:     req.ChainType,
		BlockNumber:   req.BlockNumber,
		Addresses:     req.Addresses,
		InterfaceType: req.InterfaceType,
		Status:        job.StatusCreated,
	}

	if err := h.discoveryJobsStore.CreateDiscoveryJob(discoveryJob); err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, discoveryJob)
}

func (h *Handler) ListDiscoveryJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(StatusQueryKey)
	status, err := job.ParseStatus(rawStatus)
	if err != nil {
		return parsingError(err)
	}

	jobs, err := h.discoveryJobsStore.ListDiscoveryJobs(status)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, jobs)
}

func (h *Handler) GetDiscoveryJob(ctx echo.Context) error {
	rawJobID := ctx.Param(DiscoveryJobIDParamKey)
	jobID := job.ID(rawJobID)

	discoveryJob, err := h.discoveryJobsStore.GetDiscoveryJob(jobID)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, discoveryJob)
}

func (h *Handler) CancelDiscoveryJob(ctx echo.Context) error {
	rawJobID := ctx.Param(DiscoveryJobIDParamKey)
	jobID := job.ID(rawJobID)

	if err := h.discoveryJobsStore.CancelDeliveryJob(jobID); err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) RequeueDiscoveryJob(ctx echo.Context) error {
	rawJobID := ctx.Param(DiscoveryJobIDParamKey)
	jobID := job.ID(rawJobID)

	discoveryJob, err := h.discoveryJobsStore.GetDiscoveryJob(jobID)
	if err != nil {
		return apiError(err)
	}

	discoveryJob.ID = uuid.New().String()
	discoveryJob.Status = job.StatusCreated

	if err := h.discoveryJobsStore.CreateDiscoveryJob(discoveryJob); err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, discoveryJob)
}

func (h *Handler) CreateParsingJob(ctx echo.Context) error {
	var req request.Parsing
	if err := ctx.Bind(&req); err != nil {
		return unpackError(err)
	}

	parsingJob := job.Parsing{
		ID:            uuid.New().String(),
		ChainURL:      req.ChainURL,
		ChainType:     req.ChainType,
		BlockNumber:   req.BlockNumber,
		Address:       req.Address,
		InterfaceType: req.InterfaceType,
		EventType:     req.EventType,
		Status:        job.StatusCreated,
	}

	if err := h.parsingJobsStore.CreateParsingJob(parsingJob); err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, parsingJob)
}

func (h *Handler) ListParsingJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(StatusQueryKey)
	status, err := job.ParseStatus(rawStatus)
	if err != nil {
		return parsingError(err)
	}

	jobs, err := h.parsingJobsStore.ListParsingJobs(status)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, jobs)
}

func (h *Handler) GetParsingJob(ctx echo.Context) error {
	rawJobID := ctx.Param(ParsingJobIDParamKey)
	jobID := job.ID(rawJobID)

	parsingJob, err := h.parsingJobsStore.GetParsingJob(jobID)
	if err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, parsingJob)
}

func (h *Handler) CancelParsingJob(ctx echo.Context) error {
	rawJobID := ctx.Param(ParsingJobIDParamKey)
	jobID := job.ID(rawJobID)

	if err := h.parsingJobsStore.CancelParsingJob(jobID); err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) RequeueParsingJob(ctx echo.Context) error {
	rawJobID := ctx.Param(ParsingJobIDParamKey)
	jobID := job.ID(rawJobID)

	parsingJob, err := h.parsingJobsStore.GetParsingJob(jobID)
	if err != nil {
		return apiError(err)
	}

	parsingJob.ID = uuid.New().String()
	parsingJob.Status = job.StatusCreated

	if err := h.parsingJobsStore.CreateParsingJob(parsingJob); err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, parsingJob)
}

func (h *Handler) Noop(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, nil)
}
