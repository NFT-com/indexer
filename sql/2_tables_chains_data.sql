\connect chains

CREATE TABLE IF NOT EXISTS chains
(
    id          UUID PRIMARY KEY,
    chain_id    VARCHAR(128) NOT NULL,
    name        TEXT         NOT NULL,
    description TEXT         NOT NULL,
    symbol      VARCHAR(16)  NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS marketplaces
(
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL,
    website     TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS chains_marketplaces
(
    marketplace_id UUID         NOT NULL REFERENCES marketplaces ON DELETE CASCADE,
    chain_id       UUID         NOT NULL REFERENCES chains ON DELETE CASCADE,
    address        VARCHAR(128) NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP,
    PRIMARY KEY (marketplace_id, chain_id)
);

CREATE TABLE IF NOT EXISTS collections
(
    id                     UUID PRIMARY KEY,
    chain_id               UUID         NOT NULL REFERENCES chains ON DELETE CASCADE,
    contract_collection_id VARCHAR(128) NULL,
    address                VARCHAR(128) NOT NULL,
    name                   TEXT         NOT NULL,
    description            TEXT         NOT NULL,
    symbol                 VARCHAR(16)  NOT NULL,
    slug                   VARCHAR(256) NOT NULL,
    uri                    TEXT         NOT NULL,
    website                TEXT         NOT NULL,
    image_url              TEXT         NOT NULL,
    created_at             TIMESTAMP DEFAULT NOW(),
    updated_at             TIMESTAMP
);

CREATE TABLE IF NOT EXISTS standards
(
    standard     UUID NOT NULL REFERENCES standards ON DELETE CASCADE,
    collection   UUID NOT NULL REFERENCES collections ON DELETE CASCADE,
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP,
    PRIMARY KEY(standard, collection)
);

CREATE TABLE IF NOT EXISTS standards
(
    id         UUID PRIMARY KEY,
    name       TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS event_types
(
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    standard   UUID NOT NULL REFERENCES standards ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS marketplaces_collections
(
    marketplace_id UUID NOT NULL REFERENCES marketplaces ON DELETE CASCADE,
    collection_id  UUID NOT NULL REFERENCES collections ON DELETE CASCADE,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP,
    PRIMARY KEY (marketplace_id, collection_id)
);

CREATE TABLE IF NOT EXISTS nfts
(
    id          UUID PRIMARY KEY,
    collection  UUID         NOT NULL REFERENCES collections ON DELETE CASCADE,
    token_id    VARCHAR(128) NOT NULL,
    name        TEXT         NOT NULL,
    image       TEXT         NOT NULL,
    description TEXT         NOT NULL,
    owner       VARCHAR(128) NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP,
    UNIQUE (collection, token_id)
);

CREATE TABLE IF NOT EXISTS traits
(
    id         UUID PRIMARY KEY,
    key        TEXT NOT NULL,
    value      TEXT NOT NULL,
    nft        UUID NOT NULL REFERENCES nfts ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);
