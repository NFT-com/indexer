INSERT INTO standards (id, name)
VALUES ('f7d4c503-3a75-49c8-b72b-e18b30e14d6a', 'ERC721'),
       ('4c2574d1-bd73-446b-94bb-1362f03700c0', 'ERC1155'),
       ('3f868d69-b947-4116-8104-4d984ff59756', 'OpenSea');

INSERT INTO events (id, event_hash, name)
VALUES ('fa937621-f7f9-4352-ba99-4ece180da30c', '0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef', 'Transfer(address,address,uint64)'),
       ('2f967290-0642-452c-9737-7dcb3a755064', '0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62', 'TransferSingle(address,address,address,uint256,uint256)'),
       ('dfba757d-e8d7-44b0-8102-cb107fbe4052', '0xc4109843e0b7d514e4c093114b863f8e7d8d9a458c372cd51bfe526b588006c9', 'OrdersMatched(buyHash,bytes32,address,address,uint256,bytes32)');

INSERT INTO standards_events (standard_id, event_id)
VALUES ('f7d4c503-3a75-49c8-b72b-e18b30e14d6a', 'fa937621-f7f9-4352-ba99-4ece180da30c'),
       ('4c2574d1-bd73-446b-94bb-1362f03700c0', '2f967290-0642-452c-9737-7dcb3a755064'),
       ('3f868d69-b947-4116-8104-4d984ff59756', 'dfba757d-e8d7-44b0-8102-cb107fbe4052');
