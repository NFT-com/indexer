package jobs

import (
	"github.com/NFT-com/indexer/events"
)

type Parsing struct {
	ID           string `json:"id"`
	ChainURL     string `json:"chain_url"`
	ChainType    string `json:"chain_type"`
	BlockNumber  string `json:"block_number"`
	Address      string `json:"address"`
	StandardType string `json:"standard_type"`
	EventType    string `json:"event_type"`
	Status       Status `json:"status"`
}

type ParsingResult struct {
	RawEvents    []events.RawEvent
	ParsedEvents []events.Event
}
