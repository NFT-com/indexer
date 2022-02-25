package queue

type DiscoveryJob struct {
	ChainURL   string
	ChainType  string
	StartIndex string
	EndIndex   string
	Contracts  []string
}

type ParseJob struct {
	ID              string
	NetworkID       string
	ChainID         string
	Block           uint64
	TransactionHash string
	Address         string
	Type            string
	Data            map[string]interface{}
}
