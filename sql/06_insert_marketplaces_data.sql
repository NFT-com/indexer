INSERT INTO marketplaces
    (id, name, description, website)
VALUES ('df87df1d-f0a1-4e53-b2c3-77e794a76cf2', 'OpenSea',
        'Discover, collect, and sell extraordinary NFTs on the world''s first & largest NFT marketplace.',
        'https://opensea.io/');

INSERT INTO networks_marketplaces_standards
    (network_id, marketplace_id, deployment, contract_address, start_height, standard_id)
VALUES ('94c754fe-e06c-4d2b-bb76-2faa240b5bb8', 'df87df1d-f0a1-4e53-b2c3-77e794a76cf2', 'Wyvern Exchange v1',
        '0x7be8076f4ea4a4ad08075c2508e481d6c946d12b', 5774644, '3f868d69-b947-4116-8104-4d984ff59756'),
       ('94c754fe-e06c-4d2b-bb76-2faa240b5bb8', 'df87df1d-f0a1-4e53-b2c3-77e794a76cf2', 'Wyvern Exchange v2',
        '0x7f268357A8c2552623316e2562D90e642bB538E5', 14120913, '3f868d69-b947-4116-8104-4d984ff59756'),
       ('94c754fe-e06c-4d2b-bb76-2faa240b5bb8', 'df87df1d-f0a1-4e53-b2c3-77e794a76cf2', 'Seaport 1.0',
        '0x00000000006cee72100d161c57ada5bb2be1ca79', 14801551, '78f8f51d-bbe1-4ca1-b9c3-b11abc65174d'),
       ('94c754fe-e06c-4d2b-bb76-2faa240b5bb8', 'df87df1d-f0a1-4e53-b2c3-77e794a76cf2', 'Seaport 1.1',
        '0x00000000006c3852cbEf3e08E8dF289169EdE581', 14946474, '78f8f51d-bbe1-4ca1-b9c3-b11abc65174d');
