package jobs

var (
	ColumnsParsingJobs = []string{"id", "chain_url", "chain_id", "chain_type", "block_number", "address", "interface_type", "event_type", "status"}
	ColumnsActionJobs  = []string{"id", "chain_url", "chain_id", "chain_type", "block_number", "address", "interface_type", "event_type", "token_id", "to_address", "action_type", "status"}
)
