create table loan_contract
(
    id                     serial8     not null primary key,
    loan_offer_interest_id int8        not null references loan_package_offer_interest (id),
    symbol_id              int8        not null references symbol (id),
    investor_id            varchar(20) not null,
    account_no             varchar(50) not null,
    loan_id                int8        not null default 0,
    created_at             timestamp   not null default now(),
    updated_at             timestamp   not null default now()
);

select create_updated_at_trigger('loan_contract');
select create_notify_trigger('loan_contract');
create index loan_contract_loan_interest_id on loan_contract (loan_offer_interest_id);