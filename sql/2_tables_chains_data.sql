CREATE TABLE IF NOT EXISTS chain
(
    id          VARCHAR(128) NOT NULL PRIMARY KEY,
    name        TEXT         NOT NULL,
    description TEXT         NOT NULL,
    symbol      VARCHAR(16)  NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS marketplace
(
    chain_id    VARCHAR(128) NOT NULL,
    address     VARCHAR(128) NOT NULL,
    name        TEXT         NOT NULL,
    description TEXT         NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP,
    PRIMARY KEY (chain_id, address)
);

CREATE TABLE IF NOT EXISTS collection
(
    chain_id    VARCHAR(128) NOT NULL,
    address     VARCHAR(128) NOT NULL,
    name        TEXT         NOT NULL,
    description TEXT         NOT NULL,
    symbol      VARCHAR(16)  NOT NULL,
    slug        VARCHAR(256) NOT NULL,
    standard    VARCHAR(128) NOT NULL,
    website     TEXT         NOT NULL,
    image_url   TEXT         NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP,
    PRIMARY KEY (chain_id, address, slug)
);

CREATE TABLE IF NOT EXISTS nft
(
    id         VARCHAR(128) PRIMARY KEY,
    chain_id   VARCHAR(128) NOT NULL,
    collection VARCHAR(128) NOT NULL,
    owner      VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);
