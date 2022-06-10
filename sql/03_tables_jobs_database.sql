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

CREATE INDEX parsings_chain_id_idx ON parsings(chain_id);

CREATE INDEX parsings_contract_addresses_idx ON parsings(contract_addresses);

CREATE INDEX parsings_event_hashes_idx ON parsings(event_hashes);

CREATE INDEX parsings_job_status_idx ON parsings(job_status);

CREATE INDEX parsings_start_height_idx ON parsings(start_height);

CREATE INDEX parsings_end_height_idx ON parsings(end_height);

CREATE TABLE IF NOT EXISTS actions
(
    id               UUID PRIMARY KEY,
    chain_id         NUMERIC      NOT NULL,
    contract_address VARCHAR(128) NOT NULL,
    token_id         VARCHAR(256) NOT NULL,
    action_type      VARCHAR(256) NOT NULL,
    block_height     NUMERIC      NOT NULL,
    job_status       VARCHAR(64)  NOT NULL,
    input_data       BYTEA        NOT NULL,
    status_message   TEXT,
    created_at       TIMESTAMP DEFAULT NOW(),
    updated_at       TIMESTAMP
);

CREATE INDEX actions_chain_id_idx ON actions(chain_id);

CREATE INDEX actions_contract_address_idx ON actions(contract_address);

CREATE INDEX actions_job_status_idx ON actions(job_status);

CREATE INDEX actions_block_height_idx ON actions(block_height);