/*
    traits_collections view represents trait data, expanded 
    with ID of the collection that the relevant NFT belongs to.
*/
CREATE OR REPLACE VIEW traits_collections AS
    SELECT t.id,
           t.nft_id,
           t.name,
           t.type,
           t.value,
           n.collection_id
    FROM traits t, nfts n
    WHERE t.nft_id = n.id;
