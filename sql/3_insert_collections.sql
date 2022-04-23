\connect chains

INSERT INTO collections (id, chain_id, contract_collection_id, address, name, description, symbol, slug,
                         website, image_url)
VALUES ('abc44f5b-e4c7-46c2-9b2f-629c5bd763a6', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
        '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d', 'Bored Ape Yacht Club',
        'The Bored Ape Yacht Club is a collection of 10,000 unique Bored Ape NFTs— unique digital collectibles living on the Ethereum blockchain. Your Bored Ape doubles as your Yacht Club membership card, and grants access to members-only benefits, the first of which is access to THE BATHROOM, a collaborative graffiti board. Future areas and perks can be unlocked by the community through roadmap activation. Visit www.BoredApeYachtClub.com for more details.',
        'BAYC', 'boredapeyachtclub', 'https://ipfs.io/ipfs/QmeSjSinHpPnmXmspMjwiXyN6zS4E9zccariGR3jxcaWtq/{{ .nft_id }}',
        'https://www.boredapeyachtclub.com', '');


INSERT INTO standards (id, name)
VALUES ('f7d4c503-3a75-49c8-b72b-e18b30e14d6a', 'ERC721');

INSERT INTO event_types (id, name, standard)
VALUES ('0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef', 'Transfer(address,address,uint64)',
        'f7d4c503-3a75-49c8-b72b-e18b30e14d6a');

INSERT INTO standards_collections (standard, collection)
VALUES ('f7d4c503-3a75-49c8-b72b-e18b30e14d6a','abc44f5b-e4c7-46c2-9b2f-629c5bd763a6');
