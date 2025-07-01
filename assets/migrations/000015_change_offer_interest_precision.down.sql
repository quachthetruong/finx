alter table public.loan_package_offer_interest
    alter column fee_rate type numeric(5, 4) using fee_rate::numeric(5, 4);