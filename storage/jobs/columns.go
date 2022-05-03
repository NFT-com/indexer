package jobs

var (
	ColumnsParsingJobs = []string{"id", "chain_id", "address", "event_type", "start_height", "end_height", "data", "status"}
	ColumnsActionJobs  = []string{"id", "chain_id", "address", "token_id", "height", "action_type", "data", "status"}
)
