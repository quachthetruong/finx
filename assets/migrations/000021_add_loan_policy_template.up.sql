create table loan_policy_template
(
    id                         serial8       not null primary key,
    created_at                 timestamp     not null default now(),
    updated_at                 timestamp     not null default now(),
    updated_by                 varchar(50)   not null,
    name                       varchar(80)   not null,
    interest_rate              numeric(4,4)  not null,
    interest_basis             int2          not null default 365,
    term                       int4          not null,
    pool_id_ref                int8          not null,
    overdue_interest           numeric(4,4)  not null,
    allow_extend_loan_term     boolean       not null default true,
    allow_early_payment        boolean       not null default true,
    preferential_period        int4          not null default 0,
    preferential_interest_rate numeric(4,4)  not null default 0
);

create table investor_account
(
    account_no    varchar(50) not null primary key,
    investor_id   varchar(20) not null,
    margin_status varchar(10) not null default '',
    created_at    timestamp   not null default now(),
    updated_at    timestamp   not null default now()
);

alter table loan_contract
    add column loan_product_id_ref int8 not null default 0;