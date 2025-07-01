package investor_account

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-jet/jet/v2/qrm"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/investor_account/repository"
	orderServiceRepository "financing-offer/internal/core/orderservice/repository"
)

type UseCase interface {
	VerifyAndUpdateInvestorAccountMarginStatus(ctx context.Context, request entity.InvestorAccount) (entity.InvestorAccount, error)
}

type useCase struct {
	repository             repository.InvestorAccountRepository
	orderServiceRepository orderServiceRepository.OrderServiceRepository
}

func NewUseCase(repository repository.InvestorAccountRepository, orderServiceRepository orderServiceRepository.OrderServiceRepository) UseCase {
	return &useCase{repository: repository, orderServiceRepository: orderServiceRepository}
}

func (u *useCase) VerifyAndUpdateInvestorAccountMarginStatus(ctx context.Context, request entity.InvestorAccount) (entity.InvestorAccount, error) {
	errorTemplate := "VerifyAndUpdateInvestorAccountMarginStatus useCase %w"
	// Get investor_account margin status from db
	account, getByAccountNoError := u.repository.GetByAccountNo(ctx, request.AccountNo)
	if getByAccountNoError != nil && !errors.Is(getByAccountNoError, qrm.ErrNoRows) {
		return entity.InvestorAccount{}, fmt.Errorf(errorTemplate, getByAccountNoError)
	}
	// If investor_account margin status is already version 3, return it
	if account.MarginStatus == entity.MarginStatusVersion3 {
		return account, nil
	}
	// If investor_account margin status is not version 3, or if investor_account not exist,
	// then call MO to check if investor_account margin status is version 3
	marginStatus, err := u.checkMarginStatus(ctx, request.AccountNo)
	if err != nil {
		return entity.InvestorAccount{}, fmt.Errorf(errorTemplate, err)
	}
	if errors.Is(getByAccountNoError, qrm.ErrNoRows) { // Save investor_account to db if not exist
		account, err = u.repository.Create(ctx, entity.InvestorAccount{
			AccountNo:    request.AccountNo,
			InvestorId:   request.InvestorId,
			MarginStatus: marginStatus,
		})
		if err != nil {
			return entity.InvestorAccount{}, fmt.Errorf(errorTemplate, err)
		}
	} else if marginStatus == entity.MarginStatusVersion3 { // If investor_account margin status is version 3, update it
		account.MarginStatus = marginStatus
		account, err = u.repository.Update(ctx, account)
		if err != nil {
			return entity.InvestorAccount{}, fmt.Errorf(errorTemplate, err)
		}
	}
	return account, nil
}

func (u *useCase) checkMarginStatus(ctx context.Context, accountNo string) (entity.MarginStatus, error) {
	errorTemplate := "checkMarginStatus useCase %w"
	loanPackages, err := u.orderServiceRepository.GetAllAccountLoanPackages(ctx, accountNo)
	if err != nil {
		return "", fmt.Errorf(errorTemplate, err)
	}
	// If investor_account has no loan package, return error
	if len(loanPackages) == 0 {
		return "", apperrors.ErrAccountNoInvalid
	}
	// Set investor_account margin status to version 2 by default
	marginStatus := entity.MarginStatusVersion2
	isMarginAccount := false
	for _, loanPackage := range loanPackages {
		// Check if investor_account has margin loan package
		if loanPackage.Type == "N" {
			// If loan package type is not margin, continue to next loan package
			continue
		}
		// If investor_account has margin loan package, set isMarginAccount to true
		isMarginAccount = true
		// Check if loan package has version 3 loan product
		for _, loanProduct := range loanPackage.LoanProducts {
			if _, err := getLoanProductId(loanProduct.Id); err != nil {
				// If loan product id is not number, continue to next loan product
				continue
			}
			// If loan product id is number, set investor_account margin status to version 3
			marginStatus = entity.MarginStatusVersion3
			return marginStatus, nil
		}
	}
	if !isMarginAccount {
		return "", apperrors.ErrNonMarginAccount
	}
	return marginStatus, nil
}

func getLoanProductId(accountLoanPackageId string) (int64, error) {
	results := strings.Split(accountLoanPackageId, "-")
	if len(results) < 2 {
		return 0, fmt.Errorf("invalid loan product id: %s", accountLoanPackageId)
	}
	return strconv.ParseInt(results[1], 10, 64)
}
