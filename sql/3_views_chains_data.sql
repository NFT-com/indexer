/*
    traits_collections view represents trait data, expanded 
    with ID of the collection that the relevant NFT belongs to.
*/
CREATE OR REPLACE VIEW traits_collections AS
    SELECT t.id,
            t.nft,
            t.name,
            t.value,
            n.collection
    FROM traits t, nfts n
    WHERE t.nft = n.id;
