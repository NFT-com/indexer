package postgres

const (
	DiscoveryJobsDBName = "discovery_jobs"
	ParsingJobsDBName   = "parsing_jobs"

	EventsDBName                 = "events"
	EventsTableIndexedDataColumn = "indexed_Data"
	EventsTableDataColumn        = "data"

	NFTsDBName = "nfts"
)

var (
	DiscoveryJobsTableColumns = []string{"id", "chain_url", "chain_type", "block_number", "addresses", "interface_type", "status"}
	ParsingJobsTableColumns   = []string{"id", "chain_url", "chain_type", "block_number", "address", "interface_type", "event_type", "status"}
	EventsTableColumns        = []string{"id", "chain_id", "network_id", "block_number", "block_hash", "address", "transaction_hash", "event_type"}
	NFTsTableColumns          = []string{"id", "chain_id", "network_id", "contract", "owner"}
)
