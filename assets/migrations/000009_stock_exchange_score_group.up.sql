ALTER TABLE stock_exchange
    ADD COLUMN score_group_id bigint NULL;
ALTER TABLE stock_exchange ADD CONSTRAINT fk_stock_exchange_score_group FOREIGN KEY (score_group_id) REFERENCES score_group(id);
CREATE INDEX stock_exchange_score_group_id_idx ON stock_exchange (score_group_id);

