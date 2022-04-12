\connect chains

INSERT INTO collections (id, chain_id, contract_collection_id, address, name, description, symbol, slug, uri,
                         website, image_url)
VALUES ('0556dfb2-0281-4483-b36d-0708c50593a8', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
        '0xf75dafa52b05aa5a6bad98d069fa38e38873f442', 'BUNTS', '', 'BNT', 'bunts', '', '', ''),

       ('52cfde87-2433-42dc-a5e5-7bf3d8f9933b', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
        '0x4a537F61ef574153664c0Dbc8c8F4B900cacBE5d', 'Mavia Land',
        'Heroes of Mavia is an MMO Strategy base-builder game, with play-to-earn and NFT components at the core of the Mavia ecosystem. Land is the most valuable NFT asset type inside of the Mavia ecosystem. Land is required in order to build a base, train an army and attack rival bases.',
        'LAND', 'mavialand', 'https://be.mavia.com/api/nft/metadata/{{ .nft_id }}',
        'https://www.mavia.com/', ''),

       ('14bd2888-389a-47b1-93d2-582c695e9426', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
        '0xe3086b98ba498501491f69b4bae4ef8960f77a11', 'CryptoBears Official',
        'CryptoBears is a collection of 5000 generative NFTs.', 'CB', 'cryptobearsofficiall', '',
        'https://www.cryptobearscollection.com/',
        'https://lh3.googleusercontent.com/DryBypRjydCQ4RQ5YUgzlf_3R0KQLmipkNs-bp5i_Y7LXo9iyvKxJbjZXWP82shVH6BQaX8y763e8u3TDlyXCj99XWLvhj_s4QFTGug=s120');
