/*
    traits_ratio returns, for the specified NFT,
    the ratio of NFTs from that collection that have 
    the same value for each trait.
*/
CREATE OR REPLACE FUNCTION traits_ratio(nft_id UUID)
RETURNS TABLE (
    name TEXT,          -- trait name
    value TEXT,         -- trait value
    ratio NUMERIC       -- ratio/portion of NFTs in the collection that have the same trait name/value combination
)
STABLE
STRICT
    AS $$
    DECLARE
        trait RECORD;
        total NUMERIC;
    BEGIN

        -- Retrieve the total number of NFTs from that collection.
        SELECT COUNT(*)
        FROM nfts
        WHERE collection IN (
            SELECT collection FROM nfts WHERE id = nft_id
        ) INTO total;

        -- Perform calculation for each NFT trait.
        FOR trait IN
            SELECT *
            FROM traits_collections t
            WHERE t.nft = nft_id
        LOOP
            RETURN QUERY (
                SELECT trait.name, 
                        trait.value,
                        SUM(CASE WHEN cmp.name = trait.name AND cmp.value = trait.value THEN 1 ELSE 0 END)::NUMERIC / total
                FROM traits_collections cmp
                WHERE cmp.collection = trait.collection
            );
            
        END LOOP;
    END;
$$
LANGUAGE plpgsql;
