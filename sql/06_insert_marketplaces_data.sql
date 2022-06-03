INSERT INTO marketplaces
    (id, name, description, website)
VALUES ('df87df1d-f0a1-4e53-b2c3-77e794a76cf2', 'OpenSea V2', -- Is it V2
        'Discover, collect, and sell extraordinary NFTs on the world''s first & largest NFT marketplace.',
        'https://opensea.io/'),
       ('ee039c9d-cac3-48db-b473-a93e50522a52', 'OpenSea',
        'Discover, collect, and sell extraordinary NFTs on the world''s first & largest NFT marketplace.',
        'https://opensea.io/');

INSERT INTO networks_marketplaces
    (network_id, marketplace_id, contract_address, start_height)
VALUES ('94c754fe-e06c-4d2b-bb76-2faa240b5bb8', 'df87df1d-f0a1-4e53-b2c3-77e794a76cf2',
        '0x7f268357A8c2552623316e2562D90e642bB538E5', 14120913),
       ('16afd90e-ec6a-4425-b3da-026419b776d9', 'ee039c9d-cac3-48db-b473-a93e50522a52',
        '0x7be8076f4ea4a4ad08075c2508e481d6c946d12b', 57746440000);

INSERT INTO marketplaces_standards (marketplace_id, standard_id)
VALUES ('df87df1d-f0a1-4e53-b2c3-77e794a76cf2', '3f868d69-b947-4116-8104-4d984ff59756'),
       ('ee039c9d-cac3-48db-b473-a93e50522a52', '3f868d69-b947-4116-8104-4d984ff59756');
