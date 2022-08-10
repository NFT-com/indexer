CREATE TABLE transfers
(
    id                 UUID PRIMARY KEY,
    chain_id           NUMERIC      NOT NULL,
    token_standard     VARCHAR(128) NOT NULL,
    collection_address VARCHAR(128) NOT NULL,
    token_id           VARCHAR(128) NOT NULL,
    block_number       NUMERIC      NOT NULL,
    transaction_hash   VARCHAR(128) NOT NULL,
    event_index        INTEGER      NOT NULL,
    sender_address     VARCHAR(128) NOT NULL,
    receiver_address   VARCHAR(128) NOT NULL,
    token_count        NUMERIC      NOT NULL,
    emitted_at         TIMESTAMP    NOT NULL,
    created_at         TIMESTAMP DEFAULT NOW(),
    UNIQUE (block_number, transaction_hash, event_index)
);

CREATE INDEX transfers_collection_address_idx ON transfers (LOWER(collection_address));

CREATE INDEX transfers_token_id_idx ON transfers (token_id);

CREATE TABLE sales
(
    id                  UUID PRIMARY KEY,
    chain_id            NUMERIC      NOT NULL,
    marketplace_address VARCHAR(128) NOT NULL,
    collection_address  VARCHAR(128) NOT NULL,
    token_id            VARCHAR(128) NOT NULL,
    block_number        NUMERIC      NOT NULL,
    transaction_hash    VARCHAR(128) NOT NULL,
    event_index         INTEGER      NOT NULL,
    seller_address      VARCHAR(128) NOT NULL,
    buyer_address       VARCHAR(128) NOT NULL,
    token_count         NUMERIC      NOT NULL,
    currency_value      NUMERIC      NOT NULL,
    currency_address    VARCHAR(128) NOT NULL,
    emitted_at          TIMESTAMP    NOT NULL,
    created_at          TIMESTAMP DEFAULT NOW(),
    UNIQUE (block_number, transaction_hash, event_index)
);

CREATE INDEX sales_marketplace_address_idx ON sales (LOWER(marketplace_address));

CREATE INDEX sales_collection_address_idx ON sales (LOWER(collection_address));

CREATE INDEX sales_currency_address_idx ON sales (LOWER(currency_address));

CREATE INDEX sales_token_id_idx ON sales (token_id);
