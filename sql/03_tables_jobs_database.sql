CREATE TABLE IF NOT EXISTS parsing_jobs
(
    id             UUID PRIMARY KEY,
    chain_id       VARCHAR(128)   NOT NULL,
    addresses      []VARCHAR(128) NOT NULL,
    event_types    []VARCHAR(256) NOT NULL,
    start_height   NUMERIC        NOT NULL,
    end_height     NUMERIC        NOT NULL,
    status         VARCHAR(64)    NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP
);

CREATE TABLE IF NOT EXISTS action_jobs
(
    id             UUID PRIMARY KEY,
    chain_id       VARCHAR(128) NOT NULL,
    address        VARCHAR(128) NOT NULL,
    token_id       VARCHAR(256) NOT NULL,
    action_type    VARCHAR(256) NOT NULL,
    height         NUMERIC      NOT NULL,
    status         VARCHAR(64)  NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP
);
