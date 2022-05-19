-- Collections scraped from Defillama
INSERT INTO collections
  (id, network_id, contract_address, start_height, name, description, symbol, slug, website, image_url)
VALUES
  (
    '612ecc22-36ef-4ef7-bb0b-5b864b85d089',
    '94c754fe-e06c-4d2b-bb76-2faa240b5bb8',
    '0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d',
    122875079,
    'Bored Ape Yacht Club',
    'The Bored Ape Yacht Club is a collection of 10,000 unique Bored Ape NFTs— unique digital collectibles living on the Ethereum blockchain. Your Bored Ape doubles as your Yacht Club membership card, and grants access to members-only benefits, the first of which is access to THE BATHROOM, a collaborative graffiti board. Future areas and perks can be unlocked by the community through roadmap activation. Visit www.BoredApeYachtClub.com for more details.',
    'BAYC',
    'boredapeyachtclub',
    'http://www.boredapeyachtclub.com/',
    'https://ik.imagekit.io/bayc/assets/bayc-logo.png'
  );


-- ERC-721 contracts scraped from Etherscan.io
INSERT INTO collections_standards
  (collection_id, standard_id)
VALUES
  ('612ecc22-36ef-4ef7-bb0b-5b864b85d089', 'f7d4c503-3a75-49c8-b72b-e18b30e14d6a');
