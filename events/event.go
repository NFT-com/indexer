package events

const (
	EventTypeMint   = "mint"
	EventTypeUpdate = "update"
	EventTypeBurn   = "burn"
)

type Event struct {
	Type        string `json:"type"`
	ChainID     string `json:"chain_id"`
	NftID       string `json:"nft_id"`
	Contract    string `json:"contract"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
}

type RawEvent struct {
	ID              string   `json:"id"`
	ChainID         string   `json:"chain_id"`
	BlockNumber     string   `json:"block_number"`
	BlockHash       string   `json:"block_hash"`
	Address         string   `json:"address"`
	TransactionHash string   `json:"transaction_hash"`
	EventType       string   `json:"event_type"`
	IndexData       []string `json:"index_data"`
	Data            []byte   `json:"data"`
}
