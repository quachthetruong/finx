package scheduler

import (
	"context"
	"log/slog"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/loanoffer"
)

type LoanOfferScheduler struct {
	logger       *slog.Logger
	useCase      loanoffer.UseCase
	errorService apperrors.Service
}

func NewLoanOfferScheduler(logger *slog.Logger, useCase loanoffer.UseCase, errorService apperrors.Service) *LoanOfferScheduler {
	return &LoanOfferScheduler{
		logger:       logger,
		useCase:      useCase,
		errorService: errorService,
	}
}

func (s *LoanOfferScheduler) ExpireLoanOffers() {
	if err := s.useCase.ExpireLoanOffers(context.Background()); err != nil {
		s.logger.Error("ExpireLoanOffers", slog.String("error", err.Error()))
		if err := s.errorService.NotifyError(context.Background(), err); err != nil {
			s.logger.Error("ExpireLoanOffers NotifyError", slog.String("error", err.Error()))
		}
	}
}
