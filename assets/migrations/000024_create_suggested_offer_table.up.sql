create table suggested_offer
(
    id         serial8     not null primary key,
    config_id  int8        not null references suggested_offer_config (id),
    account_no varchar(50) not null,
    symbols    jsonb       not null default '[]',
    created_at timestamp   not null default now(),
    updated_at timestamp   not null default now()
);

select create_updated_at_trigger('suggested_offer');
