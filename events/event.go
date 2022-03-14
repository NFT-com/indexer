package events

const (
	EventTypeMint   = "mint"
	EventTypeUpdate = "update"
	EventTypeBurn   = "burn"
)

type Event interface {
	Type() string
}

type RawEvent struct {
	ID              string   `json:"id"`
	ChainID         string   `json:"chain_id"`
	NetworkID       string   `json:"network_id"`
	BlockNumber     string   `json:"block_number"`
	BlockHash       string   `json:"block_hash"`
	Address         string   `json:"address"`
	TransactionHash string   `json:"transaction_hash"`
	EventType       string   `json:"event_type"`
	IndexData       []string `json:"index_data"`
	Data            []byte   `json:"data"`
}
