INSERT INTO marketplaces
    (id, name, description, website)
VALUES ('df87df1d-f0a1-4e53-b2c3-77e794a76cf2', 'OpenSea',
        'Discover, collect, and sell extraordinary NFTs on the world''s first & largest NFT marketplace.',
        'https://opensea.io/');

INSERT INTO networks_marketplaces
    (network_id, marketplace_id, contract_address)
VALUES ('94c754fe-e06c-4d2b-bb76-2faa240b5bb8', 'df87df1d-f0a1-4e53-b2c3-77e794a76cf2',
        '0x7f268357A8c2552623316e2562D90e642bB538E5');
