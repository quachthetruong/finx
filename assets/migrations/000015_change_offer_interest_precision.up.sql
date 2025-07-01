alter table public.loan_package_offer_interest
    alter column fee_rate type numeric(9, 8) using fee_rate::numeric(9, 8);