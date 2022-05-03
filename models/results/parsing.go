package results

import (
	"github.com/NFT-com/indexer/models/events"
)

type Parsing struct {
	Burns     []*events.Burn     `json:"burns"`
	Mints     []*events.Mint     `json:"mints"`
	Transfers []*events.Transfer `json:"transfers"`
	Sales     []*events.Sale     `json:"sales"`
}
