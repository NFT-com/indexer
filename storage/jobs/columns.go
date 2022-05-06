package jobs

var (
	ColumnsParsingJobs = []string{"id", "chain_id", "contract_addresses", "event_types", "start_height", "end_height", "status"}
	ColumnsActionJobs  = []string{"id", "chain_id", "contract_address", "token_id", "action_type", "block_height", "status"}
)
