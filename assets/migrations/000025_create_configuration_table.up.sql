create table financial_configuration
(
    id              serial8      not null primary key,
    attribute       varchar(100) not null,
    value           jsonb        not null,
    last_updated_by varchar(200) not null,
    created_at      timestamp    not null default now(),
    updated_at      timestamp    not null default now()
);

create unique index financial_configuration_attribute on financial_configuration (attribute);
select create_updated_at_trigger('financial_configuration');