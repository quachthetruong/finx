package apperrors

import (
	"fmt"

	"financing-offer/internal/core/entity"
)

func ErrParamInvalid(param string) AppError {
	return New(nil, WithCode(400_0001), WithMessage(fmt.Sprintf("invalid param: %s", param)))
}

var (
	ErrInvalidInvestorId                     = New(nil, WithCode(401_0002), WithMessage("invalid investor id"))
	ErrInvalidRequestStatus                  = New(nil, WithCode(400_0003), WithMessage("invalid request status"))
	ErrInvalidGuaranteedDuration             = New(nil, WithCode(400_0004), WithMessage("invalid guaranteed duration"))
	ErrBlacklistSymbolOverlap                = New(nil, WithCode(400_0005), WithMessage("affected time overlap"))
	ErrSymbolCodeNotFound                    = New(nil, WithCode(404_0006), WithMessage("symbol code not found"))
	ErrSymbolCodeInBlacklist                 = New(nil, WithCode(400_0007), WithMessage("symbol code in blacklist"))
	ErrInvalidLoanPackageOfferInterestStatus = New(
		nil, WithCode(400_0008), WithMessage("invalid loan package offer interest status"),
	)
	ErrInvestorNotAllowed = New(
		nil, WithCode(400_0009), WithMessage("investor not allowed to perform this action"),
	)
	ErrOfferExpired                       = New(nil, WithCode(400_0010), WithMessage("offer interest expired"))
	ErrAccountNoInvalid                   = New(nil, WithCode(400_0011), WithMessage("invalid accountNo"))
	ErrLoanPackageIdsInvalid              = New(nil, WithCode(400_0015), WithMessage("invalid loan package ids"))
	ErrorLoanPackageOfferInterestNotFound = New(
		nil, WithCode(404_0016), WithMessage("loan package offer interest not found"),
	)
	ErrorOfferInterestMismatch = New(
		nil, WithCode(400_0017), WithMessage("all offer interest must have the same offer"),
	)
	ErrorNotFoundLoanPackageOfferInterest   = New(nil, WithCode(404_0018), WithMessage("offer interest not found"))
	ErrMismatchAssetType                    = New(nil, WithCode(400_0020), WithMessage("mismatch asset type"))
	ErrNonMarginAccount                     = New(nil, WithCode(405_0021), WithMessage("non margin account"))
	ErrMissingLoanPolicyTemplate            = New(nil, WithCode(400_0022), WithMessage("missing loan policy template"))
	ErrLoanPolicyTemplateIdsInvalid         = New(nil, WithCode(400_0023), WithMessage("invalid loan policy template ids"))
	ErrInvalidActionType                    = New(nil, WithCode(400_0024), WithMessage("invalid action type"))
	ErrInvalidProposeType                   = New(nil, WithCode(400_0025), WithMessage("invalid propose type"))
	ErrInvalidLoanRateId                    = New(nil, WithCode(400_0026), WithMessage("invalid loan rate id"))
	ErrorLoanPackageOfferInterestIsCreating = New(nil, WithCode(400_0027), WithMessage("loan package is being created"))
	ErrorExistActiveSuggestedOfferConfig    = New(nil, WithCode(400_0028), WithMessage("exist active suggested offer config"))
	ErrorHONotActive                        = New(nil, WithCode(400_0029), WithMessage("HO is not active"))
	ErrorCreateAndAssignLoanPackageWorkflow = New(nil, WithCode(400_0030), WithMessage("error when trigger creating and assigning loan package workflow"))
	ErrorInvalidPoolId                      = New(nil, WithCode(400_0031), WithMessage("invalid pool id"))
	ErrorInvalidLoanPolicyTemplateId        = New(nil, WithCode(400_0032), WithMessage("invalid loan policy template id"))
	ErrorInvalidCurrentSubmissionStatus     = New(nil, WithCode(400_0033), WithMessage("invalid current submission status"))
	ErrorSubmissionIsNotTheLatest           = New(nil, WithCode(400_0034), WithMessage("submission is not the latest"))
)

func ErrLoanPackageAccountAlreadyExisted(accountNo string, loanId int64) AppError {
	return New(
		nil, WithCode(400_0012), WithMessage(
			fmt.Sprintf(
				"loan package account already existed for accountNo: %s, loanId: %d", accountNo, loanId,
			),
		),
	)
}

func ErrLoanPackageAccountNotExisted(loanId int64) AppError {
	return New(
		nil, WithCode(400_0013), WithMessage(
			fmt.Sprintf("loan package account not existed for loanId: %d", loanId),
		),
	)
}

func ErrInvalidFlowType(ft entity.FlowType) AppError {
	return New(nil, WithCode(400_0014), WithMessage(fmt.Sprintf("invalid flow type: %s", ft)))
}

func ErrInvalidInput(message string) AppError {
	return New(nil, WithCode(400_0019), WithMessage(message))
}
