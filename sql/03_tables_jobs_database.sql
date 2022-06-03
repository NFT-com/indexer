CREATE TABLE IF NOT EXISTS parsings
(
    id                 UUID PRIMARY KEY,
    chain_id           NUMERIC        NOT NULL,
    contract_addresses VARCHAR(128)[] NOT NULL,
    event_hashes       VARCHAR(256)[] NOT NULL,
    start_height       NUMERIC        NOT NULL,
    end_height         NUMERIC        NOT NULL,
    job_status         VARCHAR(64)    NOT NULL,
    input_data         BYTEA          NOT NULL,
    status_message     TEXT,
    created_at         TIMESTAMP DEFAULT NOW(),
    updated_at         TIMESTAMP
);

CREATE TABLE IF NOT EXISTS actions
(
    id             UUID PRIMARY KEY,
    chain_id       NUMERIC      NOT NULL,
    action_type    VARCHAR(256) NOT NULL,
    block_height   NUMERIC      NOT NULL,
    job_status     VARCHAR(64)  NOT NULL,
    input_data     BYTEA        NOT NULL,
    status_message TEXT,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP
);
