\connect chains

INSERT INTO public.marketplaces(id, name, description, WEBSITE)
VALUES ('df87df1d-f0a1-4e53-b2c3-77e794a76cf2', 'OpenSea',
        'Discover, collect, and sell extraordinary NFTs on the world''s first & largest NFT marketplace.',
        'https://opensea.io/');

INSERT INTO public.chains_marketplaces(marketplace_id, chain_id, address)
VALUES ('df87df1d-f0a1-4e53-b2c3-77e794a76cf2', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8',
        '0x7f268357a8c2552623316e2562d90e642bb538e5');

INSERT INTO marketplaces_collections (marketplace_id, collection_id)
VALUES ('df87df1d-f0a1-4e53-b2c3-77e794a76cf2', '0556dfb2-0281-4483-b36d-0708c50593a8'),
       ('df87df1d-f0a1-4e53-b2c3-77e794a76cf2', '14bd2888-389a-47b1-93d2-582c695e9426');
