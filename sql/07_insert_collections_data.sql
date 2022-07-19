-- Collections scraped from Defillama
INSERT INTO collections
  (id, network_id, contract_address, start_height, name, description, symbol, slug, website, image_url)
VALUES
  (
    '612ecc22-36ef-4ef7-bb0b-5b864b85d089',
    '94c754fe-e06c-4d2b-bb76-2faa240b5bb8',
    '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d',
    12287507,
    'Bored Ape Yacht Club',
    'The Bored Ape Yacht Club is a collection of 10,000 unique Bored Ape NFTs— unique digital collectibles living on the Ethereum blockchain. Your Bored Ape doubles as your Yacht Club membership card, and grants access to members-only benefits, the first of which is access to THE BATHROOM, a collaborative graffiti board. Future areas and perks can be unlocked by the community through roadmap activation. Visit www.BoredApeYachtClub.com for more details.',
    'BAYC',
    'boredapeyachtclub',
    'http://www.boredapeyachtclub.com/',
    'https://ik.imagekit.io/bayc/assets/bayc-logo.png'
  ),
  (
    '42321e81-38a8-4a3d-8ee9-4a9b4ad68d1b',
    '94c754fe-e06c-4d2b-bb76-2faa240b5bb8',
    '0x34d85c9CDeB23FA97cb08333b511ac86E1C4E258',
    114672945,
    'Otherdeed',
    'Otherdeed is the key to claiming land in Otherside. Each have a unique blend of environment and sediment — some with resources, some home to powerful artifacts. And on a very few, a Koda roams.',
    'OTHR',
    'otherdeed',
    'https://otherside.xyz/',
    ''
  ),
  (
    '37f5eff7-e355-4d8b-9a35-8bfa4f819fef',
    '94c754fe-e06c-4d2b-bb76-2faa240b5bb8',
    '0x87E738a3d5E5345d6212D8982205A564289e6324',
    114113471,
    'Fighter',
    'The on-chain idle MMO.',
    'FIGHTER',
    'fighter',
    'https://raid.party/',
    ''
  ),
  (
    'c34f1bd8-c0d9-47d8-b4a4-6447a019a9cd',
    '94c754fe-e06c-4d2b-bb76-2faa240b5bb8',
    '0x306b1ea3ecdf94aB739F1910bbda052Ed4A9f949',
    114492070,
    'BEANZ Official',
    'BEANZ coming soon.',
    'SMTH',
    'something',
    '',
    ''
  ),
  (
    '17d5f376-954f-4167-ac7c-0007df5efa62',
    '94c754fe-e06c-4d2b-bb76-2faa240b5bb8',
    '0xbcd4f1ecff4318e7a0c791c7728f3830db506c71',
    111662151,
    '',
    'Cometh is a DeFi powered game with yield generating NFT. Get spaceships, explore the galaxy and earn tokens.',
    'SPACESHIP',
    'cometh',
    'https://cometh.io',
    ''
  );


-- ERC-721 contracts scraped from Etherscan.io
INSERT INTO collections_standards
  (collection_id, standard_id)
VALUES
  ('612ecc22-36ef-4ef7-bb0b-5b864b85d089', 'f7d4c503-3a75-49c8-b72b-e18b30e14d6a'),
  ('42321e81-38a8-4a3d-8ee9-4a9b4ad68d1b', 'f7d4c503-3a75-49c8-b72b-e18b30e14d6a'),
  ('37f5eff7-e355-4d8b-9a35-8bfa4f819fef', 'f7d4c503-3a75-49c8-b72b-e18b30e14d6a'),
  ('c34f1bd8-c0d9-47d8-b4a4-6447a019a9cd', 'f7d4c503-3a75-49c8-b72b-e18b30e14d6a'),
  ('17d5f376-954f-4167-ac7c-0007df5efa62', 'f7d4c503-3a75-49c8-b72b-e18b30e14d6a');
