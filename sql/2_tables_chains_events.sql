\connect chains

-- Creation of mints table.
CREATE TABLE mints
(
    id               VARCHAR(128) PRIMARY KEY,
    chain_id         INT          NOT NULL,
    collection       VARCHAR(128) NOT NULL,
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
    chain_id         INT          NOT NULL,
    collection       VARCHAR(128) NOT NULL,
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
    chain_id         INT          NOT NULL,
    collection       VARCHAR(128) NOT NULL,
    block            NUMERIC      NOT NULL,
    transaction_hash VARCHAR(128) NOT NULL,
    marketplace      VARCHAR(128) NOT NULL,
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
    chain_id         INT          NOT NULL,
    collection       VARCHAR(128) NOT NULL,
    block            NUMERIC      NOT NULL,
    transaction_hash VARCHAR(128) NOT NULL,
    burner           VARCHAR(128),
    emitted_at       TIMESTAMP    NOT NULL,
    created_at       TIMESTAMP DEFAULT NOW()
);
