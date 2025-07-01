create table investor
(
    investor_id  varchar(20) not null primary key,
    custody_code varchar(20) not null,
    created_at  timestamp   not null default now(),
    updated_at  timestamp   not null default now()
);

CREATE INDEX investor_custody_code ON investor (custody_code);