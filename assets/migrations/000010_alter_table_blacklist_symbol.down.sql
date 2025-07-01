alter table blacklist_symbol
alter column status type varchar(50) using status::varchar(50);

drop type BlacklistSymbolStatus;
