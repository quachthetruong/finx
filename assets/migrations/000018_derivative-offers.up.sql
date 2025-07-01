create type asset_type as enum ('UNDERLYING', 'DERIVATIVE');

alter table symbol
    add column asset_type asset_type  not null default 'UNDERLYING',
    add column status     varchar(20) not null default 'ACTIVE';

alter table loan_package_request
    add column asset_type    asset_type    not null default 'UNDERLYING',
    add column initial_rate  numeric(6, 5) not null default 0.0,
    add column contract_size int8          not null default 0;

alter table loan_package_offer_interest
    add column asset_type    asset_type    not null default 'UNDERLYING',
    add column initial_rate  numeric(6, 5) not null default 0.0,
    add column contract_size int8          not null default 0;

