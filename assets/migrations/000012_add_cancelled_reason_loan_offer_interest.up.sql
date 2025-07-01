alter table loan_package_offer_interest
    add cancelled_reason varchar(50) not null default 'UNKNOWN';

create table logged_request
(
    id          serial8      not null primary key,
    investor_id varchar(20)  not null,
    symbol_id   int8         not null,
    reason      varchar(100) not null,
    request     jsonb        not null,
    created_at  timestamp    not null default now()
)