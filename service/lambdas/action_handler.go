package lambdas

import (
	"context"
	"fmt"

	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/service/processors"
)

type ActionHandler struct {
	metadata *processors.Metadata
}

func NewActionHandler(metadata *processors.Metadata) *ActionHandler {

	a := ActionHandler{
		metadata: metadata,
	}

	return &a
}

func (a *ActionHandler) Handle(ctx context.Context, action *jobs.Action) error {

	switch action.ActionType {

	case jobs.ActionAddition:

		return a.fetchNFT(ctx, action)

	case jobs.ActionOwnerChange:

		return a.changeOwner(ctx, action)

	default:

		return fmt.Errorf("unknown action type (%s)", action.ActionType)
	}
}

func (a *ActionHandler) fetchNFT(ctx context.Context, action *jobs.Action) error {
	return nil
}

func (a *ActionHandler) changeOwner(ctx context.Context, action *jobs.Action) error {
	return nil
}
