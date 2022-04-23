\connect chains

INSERT INTO collections (id, chain_id, contract_collection_id, address, name, description, symbol, slug, website, image_url)

VALUES

('41b30793-9a2f-4b22-9e88-1c3d79a8b763', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
'0x306b1ea3ecdf94aB739F1910bbda052Ed4A9f949', 'BEANZ Official', '', '', '', '', '')

('2968ed9c-13d4-4c4a-9b74-6f3bd9a245f5', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
'0x49cF6f5d44E70224e2E23fDcdd2C053F30aDA28B', 'Clone X', '', '', '', '', '')

('38df6b41-a4dd-4769-a250-bbf85fced1b1', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
'0xba30E5F9Bb24caa003E9f2f0497Ad287FDF95623', 'Bored Ape Kennel Club', '', '', '', '', '')

('c15574cf-b9b7-4920-94ed-934019e82363', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
'0x94B6f3978B0A32f7Fa0B15243E86af1aEc23Deb5', 'Akutar', '', '', '', '', '')

('abc44f5b-e4c7-46c2-9b2f-629c5bd763a6', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
'0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d', 'Bored Ape Yacht Club', '', '', '', '', '')

('d724a8e6-4cc5-4623-9d65-12b9ce88a137', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
'0xa3AEe8BcE55BEeA1951EF834b99f3Ac60d1ABeeB', 'VeeFriends', '', '', '', '', '')

('1dcaf640-5f89-4835-8b82-a9cf0a2fd4f9', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
'0x23581767a106ae21c074b2276D25e5C3e136a68b', 'Moonbirds', '', '', '', '', '')

('f501b362-395d-458a-ad24-6e554a97c3ad', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
'0x341A1c534248966c4b6AFaD165B98DAED4B964ef', 'Murakami.Flowers', '', '', '', '', '')

('d4b769ef-7bae-4118-9ee5-a674883a5002', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
'0x60E4d786628Fea6478F785A6d7e704777c86a7c6', 'Mutant Ape Yacht Club', '', '', '', '', '')

('1505d5a2-eecb-4f50-9f60-608e1dc87cb8', '94c754fe-e06c-4d2b-bb76-2faa240b5bb8', NULL,
'0x86825dFCa7A6224cfBd2DA48e85DF2fc3Aa7C4B1', 'RTFKT', '', '', '', '', '')

;


INSERT INTO standards (id, name)
VALUES ('f7d4c503-3a75-49c8-b72b-e18b30e14d6a', 'ERC721');

INSERT INTO event_types (id, name, standard)
VALUES ('0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef', 'Transfer(address,address,uint64)',
        'f7d4c503-3a75-49c8-b72b-e18b30e14d6a');

INSERT INTO standards_collections (standard, collection)
VALUES ('f7d4c503-3a75-49c8-b72b-e18b30e14d6a','abc44f5b-e4c7-46c2-9b2f-629c5bd763a6');
