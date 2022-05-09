package graph

var (
	ColumnsNetworks = []string{"id", "chain_id", "name", "description", "symbol"}
	ColumnsNFTs     = []string{"id", "collection_id", "token_id", "name", "uri", "image", "description", "owner"}
	ColumnsTraits   = []string{"id", "nft_id", "name", "type", "value"}
	ColumnsEvents   = []string{"id", "event_hash", "name"}
)
