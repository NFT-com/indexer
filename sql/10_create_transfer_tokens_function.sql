CREATE OR REPLACE FUNCTION transfer_tokens(sender VARCHAR(128), receiver VARCHAR(128), token_id UUID, amount NUMERIC)
    RETURNS void
    RETURNS NULL ON NULL INPUT
AS
$$

INSERT INTO owners (owner, nft_id, number)
VALUES (receiver, token_id, 0),
       (sender, token_id, 0)
ON CONFLICT DO NOTHING;


UPDATE owners
SET number =
        CASE owner
            WHEN receiver THEN number + amount
            WHEN sender THEN number - amount
            ELSE number
            END
WHERE owner IN (sender, receiver)
  AND nft_id = token_id;

$$
    LANGUAGE SQL;
