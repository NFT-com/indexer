package jobs

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/olahol/melody.v1"

	"github.com/NFT-com/indexer/service/api"
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
	validator api.Validator
}

// NewHandler returns a new API handler.
func NewHandler(wsHandler *melody.Melody, jobs JobsHandler, validator api.Validator) *Handler {
	s := Handler{
		wsHandler: wsHandler,
		jobs:      jobs,
		validator: validator,
	}

	return &s
}

// RegisterEndpoints applies the routes to the echo server.
func (h *Handler) RegisterEndpoints(server *echo.Echo) {
	websocket := server.Group("/ws")
	websocket.GET("/discoveries", h.DiscoveryWebsocketConnection)
	websocket.GET("/parsings", h.ParsingWebsocketConnection)
	websocket.GET("/additions", h.AdditionWebsocketConnection)

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
