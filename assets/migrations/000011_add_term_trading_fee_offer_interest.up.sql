alter table loan_package_offer_interest
    add term int4 not null default 0,
    add fee_rate numeric(5, 4)  not null default 0;

