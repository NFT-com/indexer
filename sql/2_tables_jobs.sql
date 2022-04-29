CREATE TABLE IF NOT EXISTS discovery_jobs
(
    id             UUID PRIMARY KEY,
    chain_url      TEXT           NOT NULL,
    chain_id       VARCHAR(128)   NOT NULL,
    chain_type     VARCHAR(256)   NOT NULL,
    block_number   NUMERIC        NOT NULL,
    addresses      VARCHAR(256)[] NOT NULL,
    interface_type VARCHAR(256)   NOT NULL,
    status         VARCHAR(64)    NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP,
    deleted_at     TIMESTAMP
);

CREATE TABLE IF NOT EXISTS parsing_jobs
(
    id             UUID PRIMARY KEY,
    chain_url      TEXT         NOT NULL,
    chain_id       VARCHAR(128) NOT NULL,
    chain_type     VARCHAR(256) NOT NULL,
    block_number   NUMERIC      NOT NULL,
    address        VARCHAR(128) NOT NULL,
    interface_type VARCHAR(256) NOT NULL,
    event_type     VARCHAR(256) NOT NULL,
    status         VARCHAR(64)  NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP,
    deleted_at     TIMESTAMP
);

CREATE TABLE IF NOT EXISTS action_jobs
(
    id             UUID PRIMARY KEY,
    chain_url      TEXT         NOT NULL,
    chain_id       VARCHAR(128) NOT NULL,
    chain_type     VARCHAR(256) NOT NULL,
    block_number   NUMERIC      NOT NULL,
    address        VARCHAR(128) NOT NULL,
    interface_type VARCHAR(256) NOT NULL,
    event_type     VARCHAR(256) NOT NULL,
    token_id       VARCHAR(256) NOT NULL,
    to_address     VARCHAR(256) NOT NULL,
    action_type    VARCHAR(256) NOT NULL,
    status         VARCHAR(64)  NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP,
    deleted_at     TIMESTAMP
);
