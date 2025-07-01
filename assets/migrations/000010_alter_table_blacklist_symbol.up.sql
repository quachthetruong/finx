create type BlacklistSymbolStatus as enum ('ACTIVE', 'INACTIVE');

alter table blacklist_symbol
alter column status type BlacklistSymbolStatus using status::BlacklistSymbolStatus;