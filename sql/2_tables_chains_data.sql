CREATE TABLE IF NOT EXISTS networks
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

CREATE TABLE IF NOT EXISTS collections
(
    id                  UUID PRIMARY KEY,
    network_id          UUID         NOT NULL REFERENCES chains ON DELETE CASCADE,
    contract_address    VARCHAR(128) NOT NULL,
    name                TEXT         NOT NULL,
    description         TEXT         NOT NULL,
    symbol              VARCHAR(16)  NOT NULL,
    slug                VARCHAR(256) NOT NULL,
    website             TEXT         NOT NULL,
    image_url           TEXT         NOT NULL,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP
);

CREATE TABLE IF NOT EXISTS nfts
(
    id             UUID         PRIMARY KEY,
    collection_id  UUID         NOT NULL REFERENCES collections ON DELETE CASCADE,
    token_id       VARCHAR(128) NOT NULL,
    name           TEXT         NOT NULL,
    uri            TEXT         NOT NULL,
    image          TEXT         NULL,
    description    TEXT         NULL,
    owner          VARCHAR(128) NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP
);

CREATE TABLE IF NOT EXISTS traits
(
    id         UUID         PRIMARY KEY,
    nft_id     UUID         NOT NULL REFERENCES nfts ON DELETE CASCADE,
    name       TEXT         NOT NULL,
    type       TEXT         NOT NULL,
    value      TEXT         NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
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
    id         TEXT         PRIMARY KEY,
    name       TEXT         NOT NULL,
    event_hash VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS networks_marketplaces
(
    network_id          UUID         NOT NULL REFERENCES networks ON DELETE CASCADE,
    marketplace_id      UUID         NOT NULL REFERENCES marketplaces ON DELETE CASCADE,
    contract_address    VARCHAR(128) NOT NULL,
    created_at          TIMESTAMP DEFAULT NOW(),
    updated_at          TIMESTAMP,
    PRIMARY KEY (marketplace_id, chain_id)
);

CREATE TABLE IF NOT EXISTS marketplaces_collections
(
    marketplace_id UUID NOT NULL REFERENCES marketplaces ON DELETE CASCADE,
    collection_id  UUID NOT NULL REFERENCES collections ON DELETE CASCADE,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP,
    PRIMARY KEY (marketplace_id, collection_id)
);

CREATE TABLE IF NOT EXISTS standards_collections
(
    standard_id     UUID NOT NULL REFERENCES standards ON DELETE CASCADE,
    collection_id   UUID NOT NULL REFERENCES collections ON DELETE CASCADE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP,
    PRIMARY KEY(standard, collection)
);

CREATE TABLE IF NOT EXISTS standards_event_types
(
    standard_id     UUID    NOT NULL REFERENCES standards ON DELETE CASCADE,
    event_type_id   TEXT    NOT NULL REFERENCES event_types ON DELETE CASCADE
);
