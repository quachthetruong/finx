create table submission_sheet_metadata
(
    id                      serial8     not null primary key,
    loan_package_request_id int8        not null references loan_package_request (id),
    status                  varchar(50) not null,
    creator                 varchar(50) not null,
    action_type             varchar(50) not null,
    flow_type               varchar(50) not null,
    propose_type            varchar(50) not null,
    created_at              timestamp   not null default now(),
    updated_at              timestamp   not null default now()
);

create index submission_sheet_metadata_loan_package_request_id_index on submission_sheet_metadata (loan_package_request_id);


create table submission_sheet_detail
(
    id                  serial8        not null primary key,
    submission_sheet_id int8           not null references submission_sheet_metadata (id),
    loan_rate           JSONB          not null,
    firm_buying_fee     numeric(6, 5)  not null,
    firm_selling_fee    numeric(6, 5)  not null,
    transfer_fee        numeric(15, 4) not null,
    loan_policies       JSONB          not null,
    comment             text           not null default '',
    created_at          timestamp      not null default now(),
    updated_at          timestamp      not null default now()
);

create index submission_sheet_detail_submission_sheet_metadata_id_index on submission_sheet_detail (submission_sheet_id);

select create_updated_at_trigger('submission_sheet_metadata');
select create_updated_at_trigger('submission_sheet_detail');

ALTER TABLE loan_package_offer_interest
    ADD COLUMN submission_sheet_detail_id INT8,
    ADD CONSTRAINT fk_submission_sheet_detail_id FOREIGN KEY (submission_sheet_detail_id) REFERENCES submission_sheet_detail(id);

create index loan_package_offer_interest_submission_sheet_detail_id_index on loan_package_offer_interest (submission_sheet_detail_id);

select audit.audit_table('submission_sheet_metadata');
select audit.audit_table('submission_sheet_detail');