package postgres

import (
	"encoding/json"
	"fmt"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/database/dbmodels/finoffer/public/model"
)

func MapSubmissionSheetMetadataEntityToDb(submissionSheetMetadata entity.SubmissionSheetMetadata) model.SubmissionSheetMetadata {
	return model.SubmissionSheetMetadata{
		ID:                   submissionSheetMetadata.Id,
		CreatedAt:            submissionSheetMetadata.CreatedAt,
		UpdatedAt:            submissionSheetMetadata.UpdatedAt,
		Creator:              submissionSheetMetadata.Creator,
		LoanPackageRequestID: submissionSheetMetadata.LoanPackageRequestId,
		Status:               string(submissionSheetMetadata.Status),
		ProposeType:          string(submissionSheetMetadata.ProposeType),
		ActionType:           string(submissionSheetMetadata.ActionType),
		FlowType:             string(submissionSheetMetadata.FlowType),
	}
}

func MapSubmissionSheetDetailEntityToDb(submissionSheetDetail entity.SubmissionSheetDetail) (model.SubmissionSheetDetail, error) {
	errorTemplate := "LoanRate JSON Marshal error: %w"
	loanRateJSON, err := json.Marshal(submissionSheetDetail.LoanRate)
	if err != nil {
		return model.SubmissionSheetDetail{}, fmt.Errorf(errorTemplate, err)
	}
	loanPoliciesJSON, err := json.Marshal(submissionSheetDetail.LoanPolicies)
	if err != nil {
		return model.SubmissionSheetDetail{}, fmt.Errorf(errorTemplate, err)
	}
	return model.SubmissionSheetDetail{
		ID:                submissionSheetDetail.Id,
		SubmissionSheetID: submissionSheetDetail.SubmissionSheetId,
		FirmSellingFee:    submissionSheetDetail.FirmSellingFee,
		FirmBuyingFee:     submissionSheetDetail.FirmBuyingFee,
		TransferFee:       submissionSheetDetail.TransferFee,
		LoanPolicies:      string(loanPoliciesJSON),
		LoanRate:          string(loanRateJSON),
		Comment:           submissionSheetDetail.Comment,
	}, nil
}

func MapSubmissionSheetMetadataDbToEntity(submissionSheetMetadata model.SubmissionSheetMetadata) entity.SubmissionSheetMetadata {
	return entity.SubmissionSheetMetadata{
		Id:                   submissionSheetMetadata.ID,
		LoanPackageRequestId: submissionSheetMetadata.LoanPackageRequestID,
		Creator:              submissionSheetMetadata.Creator,
		CreatedAt:            submissionSheetMetadata.CreatedAt,
		UpdatedAt:            submissionSheetMetadata.UpdatedAt,
		Status:               entity.SubmissionSheetStatus(submissionSheetMetadata.Status),
		ProposeType:          entity.ProposeType(submissionSheetMetadata.ProposeType),
		ActionType:           entity.ActionType(submissionSheetMetadata.ActionType),
		FlowType:             entity.FlowType(submissionSheetMetadata.FlowType),
	}
}

func MapSubmissionSheetDetailDbToEntity(submissionSheetDetail model.SubmissionSheetDetail) (entity.SubmissionSheetDetail, error) {
	var (
		loanPolicies  []entity.LoanPolicySnapShot
		loanRate      entity.LoanRate
		errorTemplate = "LoanRate JSON Unmarshal error: %w"
	)
	err := json.Unmarshal([]byte(submissionSheetDetail.LoanRate), &loanRate)
	if err != nil {
		return entity.SubmissionSheetDetail{}, fmt.Errorf(errorTemplate, err)
	}
	err = json.Unmarshal([]byte(submissionSheetDetail.LoanPolicies), &loanPolicies)
	if err != nil {
		return entity.SubmissionSheetDetail{}, fmt.Errorf(errorTemplate, err)
	}
	return entity.SubmissionSheetDetail{
		Id:                submissionSheetDetail.ID,
		LoanRate:          loanRate,
		SubmissionSheetId: submissionSheetDetail.SubmissionSheetID,
		FirmSellingFee:    submissionSheetDetail.FirmSellingFee,
		FirmBuyingFee:     submissionSheetDetail.FirmBuyingFee,
		TransferFee:       submissionSheetDetail.TransferFee,
		LoanPolicies:      loanPolicies,
		Comment:           submissionSheetDetail.Comment,
	}, nil
}
