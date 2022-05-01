package events

var (
	ColumnsMintEvents     = []string{"id", "block", "event_index", "transaction_hash", "collection", "token_id", "owner", "emitted_at"}
	ColumnsTransferEvents = []string{"id", "block", "event_index", "transaction_hash", "collection", "token_id", "from_address", "to_address", "emitted_at"}
	ColumnsSaleEvents     = []string{"id", "block", "event_index", "transaction_hash", "marketplace", "seller", "buyer", "price", "emitted_at"}
	ColumnsBurnEvents     = []string{"id", "block", "event_index", "transaction_hash", "collection", "token_id", "emitted_at"}
)
