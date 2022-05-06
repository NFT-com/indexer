INSERT INTO standards (id, name)
VALUES ('f7d4c503-3a75-49c8-b72b-e18b30e14d6a', 'ERC721'),
       ('f7d4c503-3a75-49c8-b72b-e18b30e14d6b', 'ERC1155');

INSERT INTO events (id, event_hash, name)
VALUES ('fa937621-f7f9-4352-ba99-4ece180da30c', '0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef', 'Transfer(address,address,uint64)'),
       ('2f967290-0642-452c-9737-7dcb3a755064', '0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62', 'TransferSingle(address,address,address,uint256,uint256)');

INSERT INTO standards_events (standard_id, event_id)
VALUES ('f7d4c503-3a75-49c8-b72b-e18b30e14d6a', 'fa937621-f7f9-4352-ba99-4ece180da30c'),
       ('f7d4c503-3a75-49c8-b72b-e18b30e14d6b', '2f967290-0642-452c-9737-7dcb3a755064');