create table suggested_offer_config
(
    id              serial8       not null primary key,
    name            varchar(50)   not null,
    value           numeric(6, 5) not null,
    value_type      varchar(20)   not null,
    status          varchar(20)   not null default 'INACTIVE',
    created_by      varchar(512)  not null,
    last_updated_by varchar(512)  not null default '',
    created_at      timestamp     not null default now(),
    updated_at      timestamp     not null default now()
);

select create_updated_at_trigger('suggested_offer_config');