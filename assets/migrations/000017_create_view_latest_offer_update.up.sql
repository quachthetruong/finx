create view latest_offer_update as
with partition_table as
         (select offline_offer_update.*,
                 ROW_NUMBER() over (partition by offer_id order by created_at desc) as rn
          from offline_offer_update)
select partition_table.*
from partition_table
where rn = 1;