ALTER TABLE  loan_package_offer_interest
DROP CONSTRAINT fk_submission_sheet_detail_id;

ALTER TABLE  loan_package_offer_interest
DROP COLUMN submission_sheet_detail_id;

drop table submission_sheet_detail;
drop table submission_sheet_metadata;