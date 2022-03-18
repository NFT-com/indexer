package postgres

const (
	UniqueViolation = "23505"

	DiscoveryJobsDBName = "discovery_jobs"
	ParsingJobsDBName   = "parsing_jobs"

	NFTsDBName = "nfts"

	HistoryDBName = "history"
)

var (
	DiscoveryJobsTableColumns = []string{"id", "chain_url", "chain_type", "block_number", "addresses", "interface_type", "status"}
	ParsingJobsTableColumns   = []string{"id", "chain_url", "chain_type", "block_number", "address", "interface_type", "event_type", "status"}
	NFTsTableColumns          = []string{"id", "chain_id", "network_id", "contract", "owner"}
	HistoryTableColumns       = []string{"id", "chain_id", "network_id", "event_type", "contract", "nft_id", "from_address", "to_address", "price", "emitted_at"}
)
