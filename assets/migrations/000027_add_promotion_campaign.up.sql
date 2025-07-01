create table promotion_campaign
(
    id                         serial8       not null primary key,
    created_at                 timestamp     not null default now(),
    updated_at                 timestamp     not null default now(),
    updated_by                 varchar(50)   not null,
    name                       varchar(100)  not null,
    tag                        varchar(10)   not null,
    description                varchar(255)  not null,
    status                     varchar(10)   not null,
    metadata                   jsonb         not null default '{}'
);

select create_updated_at_trigger('promotion_campaign');