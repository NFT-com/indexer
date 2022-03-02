CREATE TABLE IF NOT EXISTS discovery_jobs
(
    id             UUID PRIMARY KEY,
    chain_url      TEXT         NOT NULL,
    chain_type     VARCHAR(256) NOT NULL,
    block_number   VARCHAR(128) NOT NULL, -- This is  uint256 value -> it needs ~80 chars to represent the maximum value
    addresses      JSONB        NOT NULL,
    interface_type VARCHAR(256) NOT NULL,
    status         VARCHAR(64)  NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP,
    deleted_at     TIMESTAMP
);

CREATE TABLE IF NOT EXISTS parsing_jobs
(
    id             UUID PRIMARY KEY,
    chain_url      TEXT         NOT NULL,
    chain_type     VARCHAR(256) NOT NULL,
    block_number   VARCHAR(128) NOT NULL, -- This is  uint256 value -> it needs ~80 chars to represent the maximum value
    address        VARCHAR(128) NOT NULL,
    interface_type VARCHAR(256) NOT NULL,
    event_type     VARCHAR(256) NOT NULL,
    status         VARCHAR(64)  NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP,
    deleted_at     TIMESTAMP
);