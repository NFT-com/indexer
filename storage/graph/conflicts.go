package graph

const (
	ConflictNFTs   = "ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, uri = EXCLUDED.uri, image = EXCLUDED.image, description = EXCLUDED.description, owner = EXCLUDED.owner"
	ConflictTraits = "ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, value = EXCLUDED.value, nft = EXCLUDED.nft"
)
