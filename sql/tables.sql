-- Creation of chain table
CREATE TABLE IF NOT EXISTS chain (
    id UUID PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    symbol VARCHAR(16) NOT NULL,
    network_id TEXT NOT NULL,
    chain_id TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Creation of marketplace table
CREATE TABLE IF NOT EXISTS marketplace (
    id UUID PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Creation of collection table
CREATE TABLE IF NOT EXISTS collection (
    id UUID PRIMARY KEY,
    chain_id UUID NOT NULL, -- FOREIGN
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    symbol VARCHAR(16) NOT NULL,
    address VARCHAR(64) NOT NULL,
    abi TEXT NOT NULL,
    standard VARCHAR(16) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Creation of the junction table for marketplaces and collections
CREATE TABLE IF NOT EXISTS marketplace_collections (
    marketplace_id UUID NOT NULL, -- FOREIGN
    collection_id UUID NOT NULL -- FOREIGN
);

-- Creation of event table
CREATE TABLE IF NOT EXISTS event (
    id VARCHAR(64) PRIMARY KEY,
    chain_id UUID NOT NULL, -- FOREIGN
    collection_id UUID NOT NULL, -- FOREIGN
    block BIGINT NOT NULL,
    transaction_hash VARCHAR(64) NOT NULL,
    type VARCHAR(64) NOT NULL,
    data JSONB,
    emmited_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Creation of nft table
CREATE TABLE IF NOT EXISTS nft (
    id UUID PRIMARY KEY,
    collection_id UUID NOT NULL, -- FOREIGN
    token_id TEXT NOT NULL, -- ID of the nft in the collection
    owner VARCHAR(64) NOT NULL,
    name TEXT NOT NULL,
    uri TEXT NOT NULL,
    rarity DOUBLE PRECISION,
    data JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

-- FIXME: Should we use FOREIGN KEY?
-- FIXME: How do we map the owner of an nft to a nft?
-- FIXME: Should we use a nft_address table? nft_owner(id BIGSERIAL, nft_id BIGINT, owner)