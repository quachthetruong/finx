create table loan_request_scheduler_config
(
    id                serial8       not null primary key,
    maximum_loan_rate numeric(4, 3) not null,
    affected_from     timestamp     not null default now(),
    created_at        timestamp     not null default now(),
    updated_at        timestamp     not null default now()
);

insert into loan_request_scheduler_config(maximum_loan_rate)
values (0.9);