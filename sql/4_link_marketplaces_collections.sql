-- NOTE: This query should not be executed blindly but for bootstrapping
-- it's okay since all collections we're currently aware of are on the 
-- Opensea marketplace.
insert into marketplaces_collections (marketplace_id, collection_id)
select 'df87df1d-f0a1-4e53-b2c3-77e794a76cf2', coll.id
from collections coll;
