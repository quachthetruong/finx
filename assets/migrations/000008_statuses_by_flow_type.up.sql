alter table loan_package_request
    add guaranteed_duration int4 not null default 0;

alter table loan_package_offer
    add expired_at timestamp,
    add flow_type  varchar(50) not null default '';

alter table loan_contract
    add guaranteed_end_at timestamp,
    add loan_package_account_id int8 not null;

alter table loan_package_offer_interest
    add loan_id      int8         not null default 0,
    add cancelled_by varchar(500) not null default '',
    add cancelled_at timestamp;

select audit.audit_table('loan_package_offer');