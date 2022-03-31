package api

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/olahol/melody.v1"
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
	wsHandler *melody.Melody
	jobs      JobsHandler
	validator Validator
}

// NewHandler returns a new API handler.
func NewHandler(wsHandler *melody.Melody, jobs JobsHandler, validator Validator) *Handler {
	s := Handler{
		wsHandler: wsHandler,
		jobs:      jobs,
		validator: validator,
	}

	return &s
}

// ApplyRoutes applies the routes to the echo server.
func (h *Handler) ApplyRoutes(server *echo.Echo) {
	websocket := server.Group("/ws")
	websocket.GET("/discoveries", h.DiscoveryWebsocketConnection)
	websocket.GET("/parsings", h.ParsingWebsocketConnection)

	discoveries := server.Group("/discoveries")
	discoveries.POST("", h.CreateDiscoveryJob)
	discoveries.GET("", h.ListDiscoveryJobs)
	discoveries.GET("/highest", h.GetHighestBlockNumberDiscoveryJob)
	discoveries.GET("/:id", h.GetDiscoveryJob)
	discoveries.PATCH("/:id", h.UpdateDiscoveryJobStatus)

	parsings := server.Group("/parsings")
	parsings.POST("", h.CreateParsingJob)
	parsings.GET("", h.ListParsingJobs)
	parsings.GET("/highest", h.GetHighestBlockNumberParsingJob)
	parsings.GET("/:id", h.GetParsingJob)
	parsings.PATCH("/:id", h.UpdateParsingJobStatus)
}
