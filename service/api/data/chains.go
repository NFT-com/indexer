package data

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/NFT-com/indexer/models/chain"
	"github.com/NFT-com/indexer/service/api"
	"github.com/NFT-com/indexer/service/request"
)

func (h *Handler) CreateChain(ctx echo.Context) error {
	var req request.Chain
	err := ctx.Bind(&req)
	if err != nil {
		return api.BadRequest(err)
	}

	err = h.validator.Request(req)
	if err != nil {
		return api.BadRequest(err)
	}

	chain := chain.Chain{
		ChainID:     req.ChainID,
		Name:        req.Name,
		Description: req.Description,
		Symbol:      req.Symbol,
	}

	newChain, err := h.dataController.CreateChain(chain)
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.JSON(http.StatusOK, *newChain)
}

func (h *Handler) ListChains(ctx echo.Context) error {
	chains, err := h.dataController.ListChains()
	if err != nil {
		return api.InternalError(err)
	}

	return ctx.JSON(http.StatusOK, chains)
}
