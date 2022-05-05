package graph

var (
	ColumnsNetworks    = []string{"id", "chain_id", "name", "description", "symbol"}
	ColumnsCollections = []string{"id", "network_id", "base_token_id", "contract_address", "name", "description", "symbol", "slug", "image_url", "website"}
	ColumnsNFTs        = []string{"id", "token_id", "collection_id", "name", "uri", "image", "description", "owner"}
	ColumnsTraits      = []string{"id", "name", "value", "nft"}
	ColumnsEventTypes  = []string{"id", "name"}
)
