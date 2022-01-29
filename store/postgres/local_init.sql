-- Creation of chain table
CREATE TABLE IF NOT EXISTS chain (
    id INT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    symbol VARCHAR(16) NOT NULL,
    network VARCHAR(64) NOT NULL
)

-- Creation of event table
CREATE TABLE IF NOT EXISTS event (
    id VARCHAR(66) PRIMARY KEY,
    network_id INT NOT NULL, -- FOREIGN
    block NUMERIC NOT NULL,
    transaction_hash VARCHAR(66) NOT NULL,
    address VARCHAR(66) NOT NULL,
    type VARCHAR(64) NOT NULL,
    data JSONB
);

-- Creation of collection table
CREATE TABLE IF NOT EXISTS collection (
    id INT PRIMARY KEY,
)