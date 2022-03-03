package postgres

const (
	DiscoveryJobsDBName = "discovery_jobs"
	ParsingJobsDBName   = "parsing_jobs"
)

var (
	DiscoveryJobsTableColumns = []string{"id", "chain_url", "chain_type", "block_number", "addresses", "interface_type", "status"}
	ParsingJobsTableColumns   = []string{"id", "chain_url", "chain_type", "block_number", "address", "interface_type", "event_type", "status"}
)
