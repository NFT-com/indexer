package dispatch

type DiscoveryJob struct {
	ChainURL   string
	ChainType  string
	StartIndex string
	EndIndex   string
	Contracts  []string
}
