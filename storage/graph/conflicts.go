package graph

const (
	ConflictNFTs      = "ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, uri = EXCLUDED.uri, image = EXCLUDED.image, description = EXCLUDED.description"
	ConflictNFTOwners = "ON CONFLICT (nft_id, owner) DO UPDATE SET number = (nft_owners.number + EXCLUDED.number), updated_at = NOW()"
	ConflictTraits    = "ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, value = EXCLUDED.value, nft = EXCLUDED.nft"
)
