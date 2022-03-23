package postgres

const (
	uniqueViolation = "23505"

	discoveryJobsTableName = "discovery_jobs"
	parsingJobsTableName   = "parsing_jobs"

	nftsTableName = "nfts"

	historyDBName = "history"
)

var (
	discoveryJobsTableColumns = []string{"id", "chain_url", "chain_type", "block_number", "addresses", "interface_type", "status"}
	parsingJobsTableColumns   = []string{"id", "chain_url", "chain_type", "block_number", "address", "interface_type", "event_type", "status"}
	nftsTableColumns          = []string{"id", "chain_id", "network_id", "contract", "owner"}
	historyTableColumns       = []string{"id", "chain_id", "network_id", "event_type", "contract", "nft_id", "from_address", "to_address", "price", "emitted_at"}
)
