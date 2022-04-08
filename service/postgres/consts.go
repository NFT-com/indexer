package postgres

const (
	// jobs table names
	discoveryJobsTableName = "discovery_jobs"
	parsingJobsTableName   = "parsing_jobs"

	// events table names
	mintEventTableName     = "mints"
	transferEventTableName = "transfers"
	saleEventTableName     = "sales"
	burnEventTableName     = "burns"

	// data table names
	chainTableName       = "chains"
	collectionTableName  = "collections"
	marketplaceTableName = "marketplaces"
	nftTableName         = "nfts"
)

var (
	// jobs table columns
	discoveryJobsTableColumns = []string{"id", "chain_url", "chain_type", "block_number", "addresses", "interface_type", "status"}
	parsingJobsTableColumns   = []string{"id", "chain_url", "chain_type", "block_number", "address", "interface_type", "event_type", "status"}

	// events table columns
	mintEventTableColumns     = []string{"id", "block", "event_index", "transaction_hash", "collection", "token_id", "owner", "emitted_at"}
	transferEventTableColumns = []string{"id", "block", "event_index", "transaction_hash", "collection", "token_id", "from_address", "to_address", "emitted_at"}
	saleEventTableColumns     = []string{"id", "block", "event_index", "transaction_hash", "marketplace", "seller", "buyer", "price", "emitted_at"}
	burnEventTableColumns     = []string{"id", "block", "event_index", "transaction_hash", "collection", "token_id", "emitted_at"}

	// events on conflict statements
	mintTableOnConflictStatement     = "ON CONFLICT (id) DO UPDATE SET block = EXCLUDED.block, event_index = EXCLUDED.event_index, transaction_hash = EXCLUDED.transaction_hash, collection = EXCLUDED.collection, token_id = EXCLUDED.token_id, owner = EXCLUDED.owner, emitted_at = EXCLUDED.emitted_at"
	transferTableOnConflictStatement = "ON CONFLICT (id) DO UPDATE SET block = EXCLUDED.block, event_index = EXCLUDED.event_index, transaction_hash = EXCLUDED.transaction_hash, collection = EXCLUDED.collection,token_id = EXCLUDED.token_id, from_address = EXCLUDED.from_address, to_address = EXCLUDED.to_address, emitted_at = EXCLUDED.emitted_at"
	saleTableOnConflictStatement     = "ON CONFLICT (id) DO UPDATE SET block = EXCLUDED.block, event_index = EXCLUDED.event_index, transaction_hash = EXCLUDED.transaction_hash, marketplace = EXCLUDED.marketplace, seller = EXCLUDED.seller, buyer = EXCLUDED.buyer, price = EXCLUDED.price, emitted_at = EXCLUDED.emitted_at"
	burnTableOnConflictStatement     = "ON CONFLICT (id) DO UPDATE SET block = EXCLUDED.block, event_index = EXCLUDED.event_index, transaction_hash = EXCLUDED.transaction_hash, collection = EXCLUDED.collection,token_id = EXCLUDED.token_id, emitted_at = EXCLUDED.emitted_at"

	// data table columns
	chainTableColumns       = []string{"id", "chain_id", "name", "description", "symbol"}
	collectionTableColumns  = []string{"id", "chain_id", "contract_collection_id", "address", "name", "description", "symbol", "slug", "uri", "image_url", "website"}
	marketplaceTableColumns = []string{"id", "name", "description", "website"}
	nftTableColumns         = []string{"id", "collection", "token_id", "owner", "rarity"}
)
