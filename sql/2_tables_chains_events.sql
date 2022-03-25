\connect chains

-- Creation of mints table.
CREATE TABLE mints
(
    id               VARCHAR(128) PRIMARY KEY,
    collection       UUID         NOT NULL REFERENCES collections ON DELETE CASCADE,
    block            NUMERIC      NOT NULL,
    transaction_hash VARCHAR(128) NOT NULL,
    owner            VARCHAR(128),
    emitted_at       TIMESTAMP    NOT NULL,
    created_at       TIMESTAMP DEFAULT NOW()
);

-- Creation of transfers table.
CREATE TABLE transfers
(
    id               VARCHAR(128) PRIMARY KEY,
    collection       UUID         NOT NULL REFERENCES collections ON DELETE CASCADE,
    block            NUMERIC      NOT NULL,
    transaction_hash VARCHAR(128) NOT NULL,
    from_address     VARCHAR(128) NOT NULL,
    to_address       VARCHAR(128) NOT NULL,
    emitted_at       TIMESTAMP    NOT NULL,
    created_at       TIMESTAMP DEFAULT NOW()
);

-- Creation of sales table.
CREATE TABLE sales
(
    id               VARCHAR(128) PRIMARY KEY,
    marketplace      UUID         NOT NULL REFERENCES marketplaces ON DELETE CASCADE,
    block            NUMERIC      NOT NULL,
    transaction_hash VARCHAR(128) NOT NULL,
    seller           VARCHAR(128) NOT NULL,
    buyer            VARCHAR(128) NOT NULL,
    price            NUMERIC      NOT NULL,
    emitted_at       TIMESTAMP    NOT NULL,
    created_at       TIMESTAMP DEFAULT NOW()
);

-- Creation of burns table.
CREATE TABLE burns
(
    id               VARCHAR(128) PRIMARY KEY,
    collection       UUID         NOT NULL REFERENCES collections ON DELETE CASCADE,
    block            NUMERIC      NOT NULL,
    transaction_hash VARCHAR(128) NOT NULL,
    emitted_at       TIMESTAMP    NOT NULL,
    created_at       TIMESTAMP DEFAULT NOW()
);
