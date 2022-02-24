package dispatch

type DiscoveryJob struct {
	ChainURL   string
	StartIndex string
	EndIndex   string
	Contracts  []string
}
