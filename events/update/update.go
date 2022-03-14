package update

import (
	"github.com/NFT-com/indexer/events"
)

type Update struct {
	ChainID   string
	NetworkID string
	NftID     string
	Contract  string
	ToAddress string
}

func (m Update) Type() string {
	return events.EventTypeUpdate
}
