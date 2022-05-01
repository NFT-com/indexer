package events

const (
	ConflictMintEvents     = "ON CONFLICT (id) DO UPDATE SET block = EXCLUDED.block, event_index = EXCLUDED.event_index, transaction_hash = EXCLUDED.transaction_hash, collection = EXCLUDED.collection, token_id = EXCLUDED.token_id, owner = EXCLUDED.owner, emitted_at = EXCLUDED.emitted_at"
	ConflictTransferEvents = "ON CONFLICT (id) DO UPDATE SET block = EXCLUDED.block, event_index = EXCLUDED.event_index, transaction_hash = EXCLUDED.transaction_hash, collection = EXCLUDED.collection,token_id = EXCLUDED.token_id, from_address = EXCLUDED.from_address, to_address = EXCLUDED.to_address, emitted_at = EXCLUDED.emitted_at"
	ConflictSaleEvents     = "ON CONFLICT (id) DO UPDATE SET block = EXCLUDED.block, event_index = EXCLUDED.event_index, transaction_hash = EXCLUDED.transaction_hash, marketplace = EXCLUDED.marketplace, seller = EXCLUDED.seller, buyer = EXCLUDED.buyer, price = EXCLUDED.price, emitted_at = EXCLUDED.emitted_at"
	ConflictBurnEvents     = "ON CONFLICT (id) DO UPDATE SET block = EXCLUDED.block, event_index = EXCLUDED.event_index, transaction_hash = EXCLUDED.transaction_hash, collection = EXCLUDED.collection,token_id = EXCLUDED.token_id, emitted_at = EXCLUDED.emitted_at"
)
