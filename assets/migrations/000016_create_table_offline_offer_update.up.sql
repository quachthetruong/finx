create table offline_offer_update
(
    id         serial8      not null primary key,
    offer_id   int8         not null references loan_package_offer (id),
    status     varchar(100) not null,
    category   varchar(200) not null,
    note       text         not null default '',
    created_by varchar(200) not null,
    created_at timestamp    not null default now()
);

create index offline_offer_offer_id on offline_offer_update (offer_id);