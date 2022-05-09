package events

const (
	ConflictTransferEvents = "ON CONFLICT (id) DO UPDATE SET chain_id = EXCLUDED.chain_id, collection_address = EXCLUDED.collection_address, token_id = EXCLUDED.token_id, block_number = EXCLUDED.block_number, transaction_hash = EXCLUDED.transaction_hash, event_index = EXCLUDED.event_index, sender_address = EXCLUDED.sender_address, receiver_address = EXCLUDED.receiver_address, emitted_at = EXCLUDED.emitted_at"
	ConflictSaleEvents     = "ON CONFLICT (id) DO UPDATE SET chain_id = EXCLUDED.chain_id, marketplace_address = EXCLUDED.marketplace_address, token_id = EXCLUDED.token_id, block_number = EXCLUDED.block_number, transaction_hash = EXCLUDED.transaction_hash, event_index = EXCLUDED.event_index, seller_address = EXCLUDED.seller_address, buyer_address = EXCLUDED.buyer_address, trade_price = EXCLUDED.trade_price, emitted_at = EXCLUDED.emitted_at"
)
