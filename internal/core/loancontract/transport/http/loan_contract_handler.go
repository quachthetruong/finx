package http

import (
	"log/slog"

	"financing-offer/internal/core/loancontract"
	"financing-offer/internal/handler"
)

type LoanContractHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	useCase loancontract.UseCase
}

func NewLoanContractHandler(
	baseHandler handler.BaseHandler,
	logger *slog.Logger,
	useCase loancontract.UseCase,
) *LoanContractHandler {
	return &LoanContractHandler{
		BaseHandler: baseHandler,
		logger:      logger,
		useCase:     useCase,
	}
}
