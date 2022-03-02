package data

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/models/chain"
	"github.com/NFT-com/indexer/service/api"
	"github.com/NFT-com/indexer/service/request"
)

func (h *Handler) CreateCollection(ctx echo.Context) error {
	var req request.Collection
	err := ctx.Bind(&req)
	if err != nil {
		return api.BadRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return api.BadRequest(err)
	}

	collection := chain.Collection{
		ChainID:              req.ChainID,
		ContractCollectionID: req.ContractCollectionID,
		Address:              req.Address,
		Name:                 req.Name,
		Description:          req.Description,
		Symbol:               req.Symbol,
		Slug:                 req.Slug,
		URI:                  req.URI,
		ImageURL:             req.ImageURL,
		Website:              req.Website,
	}

	newCollection, err := h.dataController.CreateCollection(collection)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.JSON(http.StatusOK, *newCollection)
}

func (h *Handler) ListCollections(ctx echo.Context) error {
	collections, err := h.dataController.ListCollections()
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.JSON(http.StatusOK, collections)
}
