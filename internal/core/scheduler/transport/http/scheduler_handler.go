package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/scheduler"
	"financing-offer/internal/handler"
)

type SchedulerHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	useCase scheduler.UseCase
}

func NewSchedulerHandler(
	baseHandler handler.BaseHandler,
	logger *slog.Logger,
	useCase scheduler.UseCase,
) *SchedulerHandler {
	return &SchedulerHandler{
		BaseHandler: baseHandler,
		logger:      logger,
		useCase:     useCase,
	}
}

// GetAllLoanRequestSchedulerConfigs godoc
//
//	@Summary		Get all loan request scheduler configs
//	@Description	Get all loan request scheduler configs
//	@Tags			loan request scheduler,admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[[]entity.LoanRequestSchedulerConfig]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-request-scheduler-config [get]
func (h *SchedulerHandler) GetAllLoanRequestSchedulerConfigs(ctx *gin.Context) {
	configs, err := h.useCase.GetAllLoanRequestSchedulerConfig(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.LoanRequestSchedulerConfig]{
			Data: configs,
		},
	)
}

// GetCurrentLoanRequestSchedulerConfig godoc
//
//	@Summary		Get current loan request scheduler config
//	@Description	Get current loan request scheduler config
//	@Tags			loan request scheduler,admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[entity.LoanRequestSchedulerConfig]
//	@Failure		400	{object}	handler.ErrorResponse
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-request-scheduler-config/current [get]
func (h *SchedulerHandler) GetCurrentLoanRequestSchedulerConfig(ctx *gin.Context) {
	config, err := h.useCase.GetCurrentLoanRequestSchedulerConfig(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.LoanRequestSchedulerConfig]{
			Data: config,
		},
	)
}

// CreateLoanRequestSchedulerConfig godoc
//
//	@Summary		Create loan request scheduler config
//	@Description	Create loan request scheduler config
//	@Tags			loan request scheduler,admin
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoanPackageSchedulerConfig	true	"request"
//	@Success		201		{object}	handler.BaseResponse[entity.LoanRequestSchedulerConfig]
//	@Failure		400		{object}	handler.ErrorResponse
//	@Failure		500		{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/loan-request-scheduler-config [post]
func (h *SchedulerHandler) CreateLoanRequestSchedulerConfig(ctx *gin.Context) {
	var request LoanPackageSchedulerConfig
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.logger.Error("create loan package scheduler config", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid request")
		return
	}
	created, err := h.useCase.CreateLoanRequestSchedulerConfig(
		ctx, entity.LoanRequestSchedulerConfig{
			MaximumLoanRate: request.MaximumLoanRate,
			AffectedFrom:    request.AffectedFrom,
		},
	)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusCreated, handler.BaseResponse[entity.LoanRequestSchedulerConfig]{
			Data: created,
		},
	)
}
