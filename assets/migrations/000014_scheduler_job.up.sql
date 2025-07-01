create table scheduler_job
(
    id            serial8      not null primary key,
    job_type      varchar(100) not null,
    job_status    varchar(50)  not null,
    trigger_by    varchar(100) not null,
    tracking_data jsonb        not null,
    created_at    timestamp    not null default now(),
    updated_at    timestamp    not null default now()
);