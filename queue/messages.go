package queue

type DiscoveryJob struct {
	ChainURL   string
	ChainType  string
	StartIndex string
	EndIndex   string
	Contracts  []string
}

type ParseJob struct {
}
