package graph

var (
	ColumnsChains      = []string{"id", "chain_id", "name", "description", "symbol"}
	ColumnsCollections = []string{"id", "chain_id", "contract_collection_id", "address", "name", "description", "symbol", "slug", "image_url", "website"}
	ColumnsNFTs        = []string{"id", "token_id", "collection", "name", "uri", "image", "description", "owner"}
	ColumnsTraits      = []string{"id", "name", "value", "nft"}
	ColumnsEventTypes  = []string{"id, name"}
)
