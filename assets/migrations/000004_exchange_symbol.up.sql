create table stock_exchange
(
    id         serial8     not null primary key,
    code       varchar(50) not null,
    min_score  int4        not null default 0,
    max_score  int4        not null default 0,
    created_at timestamp   not null default now(),
    updated_at timestamp   not null default now()
);

create unique index stock_exchange_code on stock_exchange (code);
select create_updated_at_trigger('stock_exchange');


create table symbol
(
    id                serial8     not null primary key,
    stock_exchange_id int8        not null references stock_exchange (id),
    symbol            varchar(10) not null,
    created_at        timestamp   not null default now(),
    updated_at        timestamp   not null default now()
);

create unique index symbol_symbol on symbol (symbol);
select create_updated_at_trigger('symbol');

create table symbol_score
(
    id            serial8      not null primary key,
    symbol_id     int8         not null references symbol (id),
    score         int4         not null,
    affected_from timestamp    not null default now(),
    status        varchar(50)  not null,
    type          varchar(50)  not null,
    creator       varchar(100) not null,
    created_at    timestamp    not null default now(),
    updated_at    timestamp    not null default now()
);

create index symbol_score_symbol_id on symbol_score (symbol_id);
select create_updated_at_trigger('symbol_score');


create table blacklist_symbol
(
    id            serial8     not null primary key,
    symbol_id     int8        not null references symbol (id),
    affected_from timestamp   not null default now(),
    affected_to   timestamp            default null,
    status        varchar(50) not null,
    created_at    timestamp   not null default now(),
    updated_at    timestamp   not null default now()
);

create index blacklist_symbol_symbol_id on blacklist_symbol (symbol_id);
select create_updated_at_trigger('blacklist_symbol');
