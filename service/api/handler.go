package api

import (
	"github.com/labstack/echo/v4"
	"gopkg.in/olahol/melody.v1"
)

const (
	jobIDParamKey = "id"

	statusQueryKey = "status"
)

type Handler struct {
	wsHandler *melody.Melody
	jobs      JobsHandler
	validator Validator
}

func NewHandler(wsHandler *melody.Melody, jobs JobsHandler, validator Validator) *Handler {
	s := Handler{
		wsHandler: wsHandler,
		jobs:      jobs,
		validator: validator,
	}

	return &s
}

func (h *Handler) ApplyRoutes(server *echo.Echo) {
	websocket := server.Group("/ws")
	websocket.GET("/discoveries", h.NewDiscoveryWebsocketConnection)
	websocket.GET("/parsings", h.NewParsingWebsocketConnection)

	discoveries := server.Group("/discoveries")
	discoveries.POST("", h.CreateDiscoveryJob)
	discoveries.GET("", h.ListDiscoveryJobs)
	discoveries.GET("/:id", h.GetDiscoveryJob)
	discoveries.PATCH("/:id", h.UpdateDiscoveryJobStatus)
	discoveries.POST("/:id/requeue", h.RequeueDiscoveryJob)

	parsings := server.Group("/parsings")
	parsings.POST("", h.CreateParsingJob)
	parsings.GET("", h.ListParsingJobs)
	parsings.GET("/:id", h.GetParsingJob)
	parsings.PATCH("/:id", h.UpdateParsingJobStatus)
	parsings.POST("/:id/requeue", h.RequeueParsingJob)
}
