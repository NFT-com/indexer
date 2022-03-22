package postgres

const (
	discoveryJobsTableName = "discovery_jobs"
	parsingJobsTableName   = "parsing_jobs"
)

var (
	discoveryJobsTableColumns = []string{"id", "chain_url", "chain_type", "block_number", "addresses", "interface_type", "status"}
	parsingJobsTableColumns   = []string{"id", "chain_url", "chain_type", "block_number", "address", "interface_type", "event_type", "status"}
)
