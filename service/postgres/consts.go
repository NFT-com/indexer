package postgres

const (
	uniqueViolation = "23505"

	discoveryJobsTableName = "discovery_jobs"
	parsingJobsTableName   = "parsing_jobs"

	eventsTableName              = "events"
	eventsTableIndexedDataColumn = "indexed_Data"
	eventsTableDataColumn        = "data"

	nftsTableName = "nfts"
)

var (
	discoveryJobsTableColumns = []string{"id", "chain_url", "chain_type", "block_number", "addresses", "interface_type", "status"}
	parsingJobsTableColumns   = []string{"id", "chain_url", "chain_type", "block_number", "address", "interface_type", "event_type", "status"}
	eventsTableColumns        = []string{"id", "chain_id", "network_id", "block_number", "block_hash", "address", "transaction_hash", "event_type"}
	nftsTableColumns          = []string{"id", "chain_id", "network_id", "contract", "owner"}
)
