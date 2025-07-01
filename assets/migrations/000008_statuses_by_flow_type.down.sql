alter table loan_package_request
    drop column guaranteed_duration;

alter table loan_package_offer
    drop column expired_at,
    drop column flow_type;

alter table loan_contract
    drop column guaranteed_end_at,
    drop column loan_package_account_id;

alter table loan_package_offer_interest
    drop column loan_id,
    drop column cancelled_at,
    drop column cancelled_by;

drop trigger if exists audit_trigger_row on loan_package_offer;
drop trigger if exists audit_trigger_stm on loan_package_offer;