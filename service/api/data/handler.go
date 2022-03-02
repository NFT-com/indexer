package data

import (
	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/service/api"
)

type Handler struct {
	dataController DataController
	validator      api.Validator
}

func NewHandler(dataController DataController, validator api.Validator) *Handler {
	s := Handler{
		dataController: dataController,
		validator:      validator,
	}

	return &s
}

func (h *Handler) RegisterEndpoints(server *echo.Echo) {
	chainGroup := server.Group("/chains")
	chainGroup.POST("", h.CreateChain)
	chainGroup.GET("", h.ListChains)

	collectionJobGroup := server.Group("/collections")
	collectionJobGroup.POST("", h.CreateCollection)
	collectionJobGroup.GET("", h.ListCollections)
}
