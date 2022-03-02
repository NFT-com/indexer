package api

import (
	"github.com/NFT-com/indexer/job"
	"github.com/NFT-com/indexer/service/request"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	DeliveryJobIDParamKey = "delivery_id"
	ParsingJobIDParamKey  = "parsing_id"

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
		deliveryJobGroup.POST("", h.CreateDeliveryJob)
		deliveryJobGroup.GET("", h.ListDeliveryJobs)
		deliveryJobGroup.GET("/"+DeliveryJobIDParamKey, h.GetDeliveryJob)
		deliveryJobGroup.DELETE("/"+DeliveryJobIDParamKey, h.GetDeliveryJob)
		deliveryJobGroup.POST("/"+DeliveryJobIDParamKey+"/requeue", h.Noop)
	}

	parsingsJobGroup := server.Group("/parsings")
	{
		parsingsJobGroup.POST("", h.Noop)
		parsingsJobGroup.GET("", h.Noop)
		parsingsJobGroup.GET("/"+ParsingJobIDParamKey, h.Noop)
		parsingsJobGroup.DELETE("/"+ParsingJobIDParamKey, h.Noop)
		parsingsJobGroup.POST("/"+ParsingJobIDParamKey+"/requeue", h.Noop)
	}
}

func (h *Handler) CreateDeliveryJob(ctx echo.Context) error {
	var req request.Discovery
	if err := ctx.Bind(&req); err != nil {
		return unpackError(err)
	}

	discoveryJob := job.Discovery{
		ID:            req.ID,
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

func (h *Handler) ListDeliveryJobs(ctx echo.Context) error {
	rawStatus := ctx.QueryParam(StatusQueryKey)
	status, err := job.ParseStatus(rawStatus)
	if err != nil {
		return parsingError(err)
	}

	jobs, err := h.discoveryJobsStore.ListDiscoveryJobs(status)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, jobs)
}

func (h *Handler) GetDeliveryJob(ctx echo.Context) error {
	rawJobID := ctx.Param(DeliveryJobIDParamKey)
	jobID := job.ID(rawJobID)

	deliveryJob, err := h.discoveryJobsStore.GetDiscoveryJob(jobID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, deliveryJob)
}

func (h *Handler) CancelDeliveryJob(ctx echo.Context) error {
	rawJobID := ctx.Param(DeliveryJobIDParamKey)
	jobID := job.ID(rawJobID)

	if err := h.discoveryJobsStore.CancelDeliveryJob(jobID); err != nil {
		return apiError(err)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (h *Handler) Noop(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, nil)
}
