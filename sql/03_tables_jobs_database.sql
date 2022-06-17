CREATE TABLE IF NOT EXISTS latest
(
    chain_id           NUMERIC        NOT NULL,
    contract_address   VARCHAR(128)   NOT NULL,
    event_hash         VARCHAR(256)   NOT NULL,
    block_height       NUMERIC        NOT NULL,
    last_id            UUID           NOT NULL,
    created_at         TIMESTAMP      NOT NULL,
    UNIQUE(chain_id, contract_address, event_hash)
)

CREATE TABLE IF NOT EXISTS parsings
(
    id                  UUID PRIMARY KEY,
    chain_id            NUMERIC         NOT NULL,
    contract_addresses  VARCHAR(128)[]  NOT NULL,
    event_hashes        VARCHAR(256)[]  NOT NULL,
    start_height        NUMERIC         NOT NULL,
    end_height          NUMERIC         NOT NULL,
    input_data          BYTEA           NOT NULL,
    failure_message     TEXT            NOT NULL,
    created_at          TIMESTAMP       NOT NULL,
    failed_at           TIMESTAMP       NOT NULL
)

CREATE TABLE IF NOT EXISTS additions
(
    id                  UUID PRIMARY KEY,
    chain_id            NUMERIC         NOT NULL,
    contract_address    VARCHAR(128)    NOT NULL,
    token_id            VARCHAR(128)    NOT NULL,
    block_height        NUMERIC         NOT NULL,
    new_owner           VARCHAR(128)    NOT NULL,
    token_count         NUMERIC         NOT NULL,
    failure_message     TEXT            NOT NULL,
    created_at          TIMESTAMP       NOT NULL,
    failed_at           TIMESTAMP       NOT NULL
)

CREATE TABLE IF NOT EXISTS moves
(
    id                  UUID PRIMARY KEY,
    chain_id            NUMERIC         NOT NULL,
    contract_address    VARCHAR(128)    NOT NULL,
    token_id            VARCHAR(128)    NOT NULL,
    block_height        NUMERIC         NOT NULL,
    from_owner          VARCHAR(128)    NOT NULL,
    to_owner            VARCHAR(128)    NOT NULL,
    token_count         NUMERIC         NOT NULL,
    failure_message     TEXT            NOT NULL,
    created_at          TIMESTAMP       NOT NULL,
    failed_at           TIMESTAMP       NOT NULL
)