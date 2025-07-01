select audit.audit_table('symbol');

alter table symbol
    add column last_updated_by varchar(200) not null default '';