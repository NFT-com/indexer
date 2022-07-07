package jobs

import (
	"github.com/NFT-com/indexer/models/events"
)

type Completion struct {
	ID          string         `json:"id"`
	ChainID     uint64         `json:"chain_id"`
	StartHeight uint64         `json:"block_height"`
	EndHeight   uint64         `json:"end_height"`
	EventHashes []string       `json:"event_hashes"`
	Sales       []*events.Sale `json:"sales"`
}

func (c Completion) SaleIDs() []string {
	saleIDs := make([]string, 0, len(c.Sales))
	for _, sale := range c.Sales {
		saleIDs = append(saleIDs, sale.ID)
	}
	return saleIDs
}

func (c Completion) TransactionHashes() []string {
	hashes := make([]string, 0, len(c.Sales))
	for _, sale := range c.Sales {
		hashes = append(hashes, sale.TransactionHash)
	}
	return hashes
}
