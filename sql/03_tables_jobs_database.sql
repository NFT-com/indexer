CREATE TABLE boundaries
(
    chain_id         NUMERIC      NOT NULL,
    contract_address VARCHAR(128) NOT NULL,
    event_hash       VARCHAR(256) NOT NULL,
    last_height      NUMERIC      NOT NULL,
    last_id          UUID         NOT NULL,
    updated_at       TIMESTAMP    NOT NULL,
    PRIMARY KEY (chain_id, contract_address, event_hash)
);

CREATE TABLE parsing_failures
(
    id                 UUID PRIMARY KEY,
    chain_id           NUMERIC        NOT NULL,
    start_height       NUMERIC        NOT NULL,
    end_height         NUMERIC        NOT NULL,
    contract_addresses VARCHAR(128)[] NOT NULL,
    event_hashes       VARCHAR(256)[] NOT NULL,
    failure_message    TEXT           NOT NULL
);

CREATE TABLE addition_failures
(
    id               UUID PRIMARY KEY,
    chain_id         NUMERIC      NOT NULL,
    contract_address VARCHAR(128) NOT NULL,
    token_id         VARCHAR(256) NOT NULL,
    token_standard   VARCHAR(256) NOT NULL,
    failure_message  TEXT         NOT NULL
);

CREATE TABLE completion_failures
(
    id                 UUID PRIMARY KEY,
    chain_id           NUMERIC        NOT NULL,
    start_height       NUMERIC        NOT NULL,
    end_height         NUMERIC        NOT NULL,
    transaction_hashes VARCHAR(128)[] NOT NULL,
    sale_ids           UUID[]         NOT NULL,
    failure_message    TEXT           NOT NULL
);
