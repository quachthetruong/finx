drop trigger if exists audit_trigger_row on symbol;
drop trigger if exists audit_trigger_stm on symbol;

alter table symbol
    drop column last_updated_by;