create table score_group
(
    id         serial8     not null primary key,
    code       varchar(50) not null,
    min_score  int4        not null,
    max_score  int4        not null,
    created_at timestamp   not null default now(),
    updated_at timestamp   not null default now()
);

select create_updated_at_trigger('score_group');

create table score_group_interest
(
    id             serial8        not null primary key,
    limit_amount   numeric(15, 4) not null,
    loan_rate      numeric(4, 3)  not null,
    interest_rate  numeric(4, 3)  not null,
    score_group_id int8           not null references score_group (id),
    created_at     timestamp      not null default now(),
    updated_at     timestamp      not null default now()
);

select create_updated_at_trigger('score_group_interest');
create index score_group_interest_score_group_id_index on score_group_interest (score_group_id);

create table loan_package_request
(
    id           serial8        not null primary key,
    symbol_id    int8           not null references symbol (id),
    investor_id  varchar(20)    not null,
    account_no   varchar(50)    not null,
    loan_rate    numeric(4, 3)  not null,
    limit_amount numeric(15, 4) not null,
    type         varchar(20)    not null,
    status       varchar(20)    not null,
    created_at   timestamp      not null default now(),
    updated_at   timestamp      not null default now()
);

select create_updated_at_trigger('loan_package_request');
create index loan_package_request_symbol_id_index on loan_package_request (symbol_id);
create index loan_package_request_investor_id_index on loan_package_request (investor_id);

create table loan_package_offer
(
    id                      serial8     not null primary key,
    loan_package_request_id int8        not null references loan_package_request (id),
    offered_by              varchar(50) not null,
    created_at              timestamp   not null default now(),
    updated_at              timestamp   not null default now()
);

select create_updated_at_trigger('loan_package_offer');
create index loan_package_offer_loan_package_request_id_index on loan_package_offer (loan_package_request_id);

create table loan_package_offer_interest
(
    id                      serial8        not null primary key,
    loan_package_offer_id   int8           not null references loan_package_offer (id),
    score_group_interest_id int8           null references score_group_interest (id),
    limit_amount            numeric(15, 4) not null,
    loan_rate               numeric(4, 3)  not null,
    interest_rate           numeric(4, 3)  not null,
    status                  varchar(20)    not null,
    created_at              timestamp      not null default now(),
    updated_at              timestamp      not null default now()
);

select create_updated_at_trigger('loan_package_offer_interest');
create index loan_package_offer_interest_loan_package_offer_id_index on loan_package_offer_interest (loan_package_offer_id);
create index loan_package_offer_interest_score_group_interest_id_index on loan_package_offer_interest (score_group_interest_id);