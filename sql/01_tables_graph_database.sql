CREATE TABLE networks
(
    id          UUID PRIMARY KEY,
    chain_id    NUMERIC     NOT NULL UNIQUE,
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL,
    symbol      VARCHAR(16) NOT NULL
);

CREATE TABLE marketplaces
(
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT NOT NULL,
    website     TEXT NOT NULL
);

CREATE TABLE collections
(
    id               UUID PRIMARY KEY,
    network_id       UUID         NOT NULL REFERENCES networks ON DELETE CASCADE,
    contract_address VARCHAR(128) NOT NULL,
    start_height     NUMERIC      NOT NULL,
    name             TEXT         NOT NULL,
    description      TEXT         NOT NULL,
    symbol           VARCHAR(16)  NOT NULL,
    slug             VARCHAR(256) NOT NULL,
    website          TEXT         NOT NULL,
    image_url        TEXT         NOT NULL,
    UNIQUE (network_id, contract_address)
);

CREATE INDEX collections_contract_address_idx ON collections(LOWER(contract_address));

CREATE TABLE nfts
(
    id            UUID PRIMARY KEY,
    collection_id UUID         NOT NULL REFERENCES collections ON DELETE CASCADE,
    token_id      VARCHAR(128) NOT NULL,
    name          TEXT         NOT NULL,
    uri           TEXT         NOT NULL,
    image         TEXT         NOT NULL,
    description   TEXT         NOT NULL,
    created_at    TIMESTAMP,
    updated_at    TIMESTAMP,
    UNIQUE (collection_id, token_id)
);

CREATE TABLE owners
(
    owner  VARCHAR(128) NOT NULL,
    nft_id UUID         NOT NULL REFERENCES nfts ON DELETE CASCADE,
    event_id UUID       NOT NULL,
    number NUMERIC      NOT NULL,
    PRIMARY KEY (owner, nft_id, event_id)
);

CREATE INDEX owners_nft_id_idx ON owners(nft_id);

CREATE TABLE traits
(
    id     UUID PRIMARY KEY,
    nft_id UUID NOT NULL REFERENCES nfts ON DELETE CASCADE,
    name   TEXT NOT NULL,
    type   TEXT NOT NULL,
    value  TEXT NOT NULL
);

CREATE INDEX traits_nft_id_idx ON traits(nft_id);

CREATE TABLE standards
(
    id   UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE events
(
    id         UUID PRIMARY KEY,
    event_hash VARCHAR(128) NOT NULL,
    name       TEXT         NOT NULL
);

CREATE TABLE networks_marketplaces
(
    network_id       UUID         NOT NULL REFERENCES networks ON DELETE CASCADE,
    marketplace_id   UUID         NOT NULL REFERENCES marketplaces ON DELETE CASCADE,
    contract_address VARCHAR(128) NOT NULL,
    start_height     NUMERIC      NOT NULL,
    PRIMARY KEY (network_id, marketplace_id, contract_address)
);

CREATE TABLE marketplaces_standards
(
    marketplace_id UUID NOT NULL REFERENCES marketplaces ON DELETE CASCADE,
    standard_id    UUID NOT NULL REFERENCES standards ON DELETE CASCADE,
    PRIMARY KEY (marketplace_id, standard_id)
);

CREATE TABLE marketplaces_collections
(
    marketplace_id UUID NOT NULL REFERENCES marketplaces ON DELETE CASCADE,
    collection_id  UUID NOT NULL REFERENCES collections ON DELETE CASCADE,
    PRIMARY KEY (marketplace_id, collection_id)
);

CREATE TABLE collections_standards
(
    collection_id UUID NOT NULL REFERENCES collections ON DELETE CASCADE,
    standard_id   UUID NOT NULL REFERENCES standards ON DELETE CASCADE,
    PRIMARY KEY (collection_id, standard_id)
);

CREATE TABLE standards_events
(
    standard_id UUID NOT NULL REFERENCES standards ON DELETE CASCADE,
    event_id    UUID NOT NULL REFERENCES events ON DELETE CASCADE,
    PRIMARY KEY (standard_id, event_id)
);

CREATE TABLE currencies
(
    id          UUID PRIMARY KEY,
    network_id  UUID            NOT NULL REFERENCES networks ON DELETE CASCADE,
    name        TEXT            NOT NULL,
    symbol      VARCHAR(16)     NOT NULL,
    address     VARCHAR(128)    NOT NULL,
    decimals    INTEGER         NOT NULL,
    endpoint    TEXT            NOT NULL,
    UNIQUE(network_id, address)
);
