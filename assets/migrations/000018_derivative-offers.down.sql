alter table loan_package_offer_interest
    drop column asset_type,
    drop column initial_rate,
    drop column contract_size;

alter table loan_package_request
    drop column asset_type,
    drop column initial_rate,
    drop column contract_size;

alter table symbol
    drop column asset_type,
    drop column status;

drop type asset_type;
