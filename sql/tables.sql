-- Creation of chain table
CREATE TABLE IF NOT EXISTS chain (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    symbol VARCHAR(16) NOT NULL,
    network VARCHAR(64) NOT NULL
);

-- Creation of collection table
CREATE TABLE IF NOT EXISTS collection (
    id BIGSERIAL PRIMARY KEY,
    network_id BIGINT NOT NULL, -- FOREIGN
    name VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    symbol VARCHAR(5) NOT NULL,
    address VARCHAR(66) NOT NULL,
    abi TEXT NOT NULL,
    standard VARCHAR(16) NOT NULL
);

-- Creation of event table
CREATE TABLE IF NOT EXISTS event (
    id VARCHAR(66) PRIMARY KEY,
    network_id BIGINT NOT NULL, -- FOREIGN
    collection_id BIGINT NOT NULL, -- FOREIGN
    block NUMERIC NOT NULL,
    transaction_hash VARCHAR(66) NOT NULL,
    type VARCHAR(64) NOT NULL,
    data JSONB
);

-- Creation of nft table
CREATE TABLE IF NOT EXISTS nft (
    id BIGSERIAL PRIMARY KEY,
    collection_id BIGINT NOT NULL, -- FOREIGN
    nft_collection_id INT NOT NULL, -- ID of the nft in the collection
    name VARCHAR(64) NOT NULL,
    uri TEXT NOT NULL,
    data JSONB
);

-- FIXME: Should we use FOREIGN KEY?
-- FIXME: How do we map the owner of an nft to a nft?
-- FIXME: Should we use a nft_address table? nft_owner(id BIGSERIAL, nft_id BIGINT, owner)