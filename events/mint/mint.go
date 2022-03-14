package mint

import (
	"github.com/NFT-com/indexer/events"
)

type Mint struct {
	ChainID   string
	NetworkID string
	NftID     string
	Contract  string
	ToAddress string
}

func (m Mint) Type() string {
	return events.EventTypeMint
}
