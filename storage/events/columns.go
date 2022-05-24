package events

var (
	ColumnsTransferEvents = []string{
		"id",
		"chain_id",
		"collection_address",
		"token_id",
		"block_number",
		"transaction_hash",
		"event_index",
		"sender_address",
		"receiver_address",
		"emitted_at",
	}

	ColumnsSaleEvents = []string{
		"id",
		"chain_id",
		"marketplace_address",
		"collection_address",
		"token_id",
		"block_number",
		"transaction_hash",
		"event_index",
		"seller_address",
		"buyer_address",
		"trade_price",
		"trade_currency",
		"emitted_at",
	}
)
