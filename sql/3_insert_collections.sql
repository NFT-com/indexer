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

       ('df65ac20-e39c-441c-bee3-6cacfb7fa991', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
        '0x90d4ffbf13bf3203940e6dace392f7c23ff6b9ed', 'Cupcat Kittens',
        'Cupcat Kittens are a collection made by Cupcats as 2nd season. This collection includes cute kittens that are part of Cupcats ecosystem.',
        'CCK', 'cupcatkittens',
        'https://cupcat.mypinata.cloud/ipfs/QmWzju1QTCYmNU59WGYw8CoGabhKPv1SDwdKTY8ow4s3sy/{{ .nft_id }}/',
        '', ''),

       ('e0ddf773-d4d9-4749-ae2f-17dc90ced1f0', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
        '0x87E738a3d5E5345d6212D8982205A564289e6324', 'Fighter',
        'The on-chain idle MMO.',
        'FIGHTER', 'fighter', 'https://api.raid.party/metadata/fighter/{{ .nft_id }}',
        'https://raid.party/', ''),

       ('14bd2888-389a-47b1-93d2-582c695e9426', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
        '0xe3086b98ba498501491f69b4bae4ef8960f77a11', 'CryptoBears Official',
        'CryptoBears is a collection of 5000 generative NFTs.', 'CB', 'cryptobearsofficiall', '',
        'https://www.cryptobearscollection.com/',
        'https://lh3.googleusercontent.com/DryBypRjydCQ4RQ5YUgzlf_3R0KQLmipkNs-bp5i_Y7LXo9iyvKxJbjZXWP82shVH6BQaX8y763e8u3TDlyXCj99XWLvhj_s4QFTGug=s120');


INSERT INTO standards (id, name)
VALUES ('f7d4c503-3a75-49c8-b72b-e18b30e14d6a', 'ERC721');

INSERT INTO event_types (id, name, standard)
VALUES ('0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef', 'Transfer(address,address,uint64)',
        'f7d4c503-3a75-49c8-b72b-e18b30e14d6a');

INSERT INTO standards_collections (standard, collection)
VALUES ('f7d4c503-3a75-49c8-b72b-e18b30e14d6a', 'df65ac20-e39c-441c-bee3-6cacfb7fa991'),
       ('f7d4c503-3a75-49c8-b72b-e18b30e14d6a', 'e0ddf773-d4d9-4749-ae2f-17dc90ced1f0');
