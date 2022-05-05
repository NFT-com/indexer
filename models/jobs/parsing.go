package jobs

// Parsing is a job that parses an NFT's data from block data.
type Parsing struct {
	ID          string   `json:"id"`
	NetworkID   string   `json:"network_id"`
	Addresses   []string `json:"addresses"`
	EventTypes  []string `json:"event_types"`
	StartHeight uint64   `json:"start_height"`
	EndHeight   uint64   `json:"end_height"`
	Data        []byte   `json:"data"`
	Status      string   `json:"status"`
}
