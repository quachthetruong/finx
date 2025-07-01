package scheduler

import (
	"context"
	"log/slog"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/loanpackagerequest"
	"financing-offer/internal/core/scheduler"
)

type LoanRequestScheduler struct {
	logger           *slog.Logger
	useCase          loanpackagerequest.UseCase
	schedulerUseCase scheduler.UseCase
	errorService     apperrors.Service
}

func NewLoanRequestScheduler(
	logger *slog.Logger,
	schedulerUseCase scheduler.UseCase,
	useCase loanpackagerequest.UseCase,
	errorService apperrors.Service,
) *LoanRequestScheduler {
	return &LoanRequestScheduler{
		logger:           logger,
		useCase:          useCase,
		schedulerUseCase: schedulerUseCase,
		errorService:     errorService,
	}
}

func (s *LoanRequestScheduler) DeclineLoanRequests() {
	ctx := context.Background()
	loanRequestSchedulerConfig, err := s.schedulerUseCase.GetCurrentLoanRequestSchedulerConfig(ctx)
	err = s.useCase.SystemDeclineRiskLoanRequests(ctx, loanRequestSchedulerConfig.MaximumLoanRate)
	if err != nil {
		s.logger.Error("DeclineLoanRequests", slog.String("error", err.Error()))
		s.notifyError(ctx, err)
	}
}

func (s *LoanRequestScheduler) notifyError(ctx context.Context, err error) {
	err = s.errorService.NotifyError(ctx, err)
	if err != nil {
		s.logger.Error("DeclineLoanRequests NotifyError", slog.String("error", err.Error()))
		return
	}
}
