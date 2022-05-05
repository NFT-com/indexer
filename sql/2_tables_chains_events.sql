-- Creation of mints table.
CREATE TABLE mints
(
    id               VARCHAR(128) PRIMARY KEY,
    chain_id         VARCHAR(128) NOT NULL,
    contract_address VARCHAR(128) NOT NULL,
    block_number     NUMERIC      NOT NULL,
    event_index      INTEGER      NOT NULL,
    transaction_hash VARCHAR(128) NOT NULL,
    token_id         VARCHAR(128) NOT NULL,
    to_address       VARCHAR(128),
    emitted_at       TIMESTAMP    NOT NULL,
    created_at       TIMESTAMP DEFAULT NOW()
);

-- Creation of transfers table.
CREATE TABLE transfers
(
    id               VARCHAR(128) PRIMARY KEY,
    chain_id         VARCHAR(128) NOT NULL,
    contract_address VARCHAR(128) NOT NULL,
    block_number     NUMERIC      NOT NULL,
    event_index      INTEGER      NOT NULL,
    transaction_hash VARCHAR(128) NOT NULL,
    token_id         VARCHAR(128) NOT NULL,
    from_address     VARCHAR(128) NOT NULL,
    to_address       VARCHAR(128) NOT NULL,
    emitted_at       TIMESTAMP    NOT NULL,
    created_at       TIMESTAMP DEFAULT NOW()
);

-- Creation of sales table.
CREATE TABLE sales
(
    id               VARCHAR(128) PRIMARY KEY,
    chain_id         VARCHAR(128) NOT NULL,
    contract_address VARCHAR(128) NOT NULL,
    block_number     NUMERIC      NOT NULL,
    event_index      INTEGER      NOT NULL,
    transaction_hash VARCHAR(128) NOT NULL,
    seller_address   VARCHAR(128) NOT NULL,
    buyer_address    VARCHAR(128) NOT NULL,
    trade_price      NUMERIC      NOT NULL,
    emitted_at       TIMESTAMP    NOT NULL,
    created_at       TIMESTAMP DEFAULT NOW()
);

-- Creation of burns table.
CREATE TABLE burns
(
    id               VARCHAR(128) PRIMARY KEY,
    chain_id         VARCHAR(128) NOT NULL,
    contract_address VARCHAR(128) NOT NULL,
    block_number     NUMERIC      NOT NULL,
    event_index      INTEGER      NOT NULL,
    transaction_hash VARCHAR(128) NOT NULL,
    token_id         VARCHAR(128) NOT NULL,
    emitted_at       TIMESTAMP    NOT NULL,
    created_at       TIMESTAMP DEFAULT NOW()
);
