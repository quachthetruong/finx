package http

import (
	"financing-offer/internal/apperrors"
	"financing-offer/internal/core/marginoperation"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loanpolicytemplate"
	"financing-offer/internal/handler"
)

type LoanPolicyTemplateHandler struct {
	handler.BaseHandler
	Logger                 *slog.Logger
	useCase                loanpolicytemplate.UseCase
	marginOperationUseCase marginoperation.UseCase
}

func NewLoanPolicyTemplateHandler(
	baseHandler handler.BaseHandler,
	logger *slog.Logger,
	useCase loanpolicytemplate.UseCase,
	marginOperationUseCase marginoperation.UseCase,
) *LoanPolicyTemplateHandler {
	return &LoanPolicyTemplateHandler{
		BaseHandler:            baseHandler,
		Logger:                 logger,
		useCase:                useCase,
		marginOperationUseCase: marginOperationUseCase,
	}
}

func (h *LoanPolicyTemplateHandler) GetAll(ctx *gin.Context) {
	policies, err := h.useCase.GetAll(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	poolIds := make([]int64, 0)
	for _, policy := range policies {
		poolIds = append(poolIds, policy.PoolIdRef)
	}
	poolGroupMap, err := h.marginOperationUseCase.GetPoolGroupMapByPoolIds(ctx, poolIds)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	aggregateLoanPolicyTemplates := make([]entity.AggregateLoanPolicyTemplate, 0)
	for _, policy := range policies {
		group, ok := poolGroupMap[policy.PoolIdRef]
		if !ok {
			h.RenderError(ctx, apperrors.ErrorInvalidPoolId)
			return
		}
		aggregateLoanPolicyTemplates = append(aggregateLoanPolicyTemplates, policy.ToAggregateModel(group))
	}

	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.AggregateLoanPolicyTemplate]{
			Data: aggregateLoanPolicyTemplates,
		},
	)
}

func (h *LoanPolicyTemplateHandler) GetById(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.Logger.Error("get loan policy template by id", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	policy, err := h.useCase.GetById(ctx, id)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.LoanPolicyTemplate]{
			Data: policy,
		},
	)
}

func (h *LoanPolicyTemplateHandler) Create(ctx *gin.Context) {
	req := CreateLoanPackageTemplateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("create loan policy template", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	err := h.marginOperationUseCase.VerifyMarginPoolId(ctx, req.PoolIdRef)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	created, err := h.useCase.Create(ctx, req.toEntity())
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusCreated, handler.BaseResponse[entity.LoanPolicyTemplate]{
			Data: created,
		},
	)
}

func (h *LoanPolicyTemplateHandler) Update(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.Logger.Error("update loan policy template", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	newLoanPolicyTemplate := UpdateLoanPackageTemplateRequest{}
	if err := ctx.ShouldBindJSON(&newLoanPolicyTemplate); err != nil {
		h.RenderBadRequest(ctx, "invalid payload", err.Error())
		return
	}
	err = h.marginOperationUseCase.VerifyMarginPoolId(ctx, newLoanPolicyTemplate.PoolIdRef)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	updated, err := h.useCase.Update(ctx, newLoanPolicyTemplate.toEntity(id))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.LoanPolicyTemplate]{
			Data: updated,
		},
	)
}

func (h *LoanPolicyTemplateHandler) Delete(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.Logger.Error("delete loan policy template", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	if err := h.useCase.Delete(ctx, id); err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusNoContent, handler.BaseResponse[string]{
			Data: "ok",
		},
	)
}
