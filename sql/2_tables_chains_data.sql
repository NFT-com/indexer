CREATE TABLE IF NOT EXISTS chains
(
    id          UUID         NOT NULL PRIMARY KEY,
    chain_id    VARCHAR(128) NOT NULL,
    name        TEXT         NOT NULL,
    description TEXT         NOT NULL,
    symbol      VARCHAR(16)  NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS marketplaces
(
    id          UUID         NOT NULL PRIMARY KEY,
    chain_id    UUID         NOT NULL REFERENCES chains ON DELETE CASCADE,
    address     VARCHAR(128) NOT NULL,
    name        TEXT         NOT NULL,
    description TEXT         NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collections
(
    id          UUID         NOT NULL PRIMARY KEY,
    chain_id    UUID         NOT NULL REFERENCES chains ON DELETE CASCADE,
    address     VARCHAR(128) NOT NULL,
    name        TEXT         NOT NULL,
    description TEXT         NOT NULL,
    symbol      VARCHAR(16)  NOT NULL,
    slug        VARCHAR(256) NOT NULL,
    standard    VARCHAR(128) NOT NULL,
    website     TEXT         NOT NULL,
    image_url   TEXT         NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS nfts
(
    id         VARCHAR(128) PRIMARY KEY,
    chain_id   UUID         NOT NULL REFERENCES chains ON DELETE CASCADE,
    collection UUID         NOT NULL REFERENCES chains ON DELETE CASCADE,
    owner      VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);
