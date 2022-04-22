package api

import (
	"github.com/labstack/echo/v4"
)

const (
	jobIDParamKey = "id"

	statusQueryKey       = "status"
	chainURLQueryKey     = "chain_url"
	chainTypeQueryKey    = "chain_type"
	addressQueryKey      = "address"
	addressesQueryKey    = "addresses"
	standardTypeQueryKey = "standard_type"
	eventTypeQueryKey    = "event_type"
)

// Handler represents the API handler.
type Handler struct {
	jobs      JobsHandler
	validator Validator
}

// NewHandler returns a new API handler.
func NewHandler(jobs JobsHandler, validator Validator) *Handler {
	s := Handler{
		jobs:      jobs,
		validator: validator,
	}

	return &s
}

// ApplyRoutes applies the routes to the echo server.
func (h *Handler) ApplyRoutes(server *echo.Echo) {
	discoveries := server.Group("/discoveries")
	discoveries.POST("", h.CreateDiscoveryJob)
	discoveries.POST("/batch", h.CreateDiscoveryJobs)
	discoveries.GET("", h.ListDiscoveryJobs)
	discoveries.GET("/highest", h.GetHighestBlockNumberDiscoveryJob)
	discoveries.GET("/:id", h.GetDiscoveryJob)
	discoveries.PATCH("/:id", h.UpdateDiscoveryJobStatus)

	parsings := server.Group("/parsings")
	parsings.POST("", h.CreateParsingJob)
	parsings.POST("/batch", h.CreateParsingJobs)
	parsings.GET("", h.ListParsingJobs)
	parsings.GET("/highest", h.GetHighestBlockNumberParsingJob)
	parsings.GET("/:id", h.GetParsingJob)
	parsings.PATCH("/:id", h.UpdateParsingJobStatus)

	additions := server.Group("/additions")
	additions.POST("", h.CreateAdditionJob)
	additions.POST("/batch", h.CreateAdditionJobs)
	additions.GET("", h.ListAdditionJobs)
	additions.GET("/:id", h.GetAdditionJob)
	additions.PATCH("/:id", h.UpdateAdditionJobStatus)
}
