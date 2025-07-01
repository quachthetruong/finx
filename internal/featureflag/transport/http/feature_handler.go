package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/featureflag"
	"financing-offer/internal/handler"
)

type FeatureHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	useCase featureflag.UseCase
}

func NewFeatureHandler(baseHandler handler.BaseHandler, logger *slog.Logger, useCase featureflag.UseCase) *FeatureHandler {
	return &FeatureHandler{
		BaseHandler: baseHandler,
		logger:      logger,
		useCase:     useCase,
	}
}

func (h *FeatureHandler) CheckFeatureEnable(ctx *gin.Context) {
	featureName := ctx.Param("name")
	if featureName == "" {
		h.RenderBadRequest(ctx, "feature name is required")
		return
	}
	investorId, err := h.InvestorId(ctx)
	if err != nil {
		h.RenderUnauthenticated(ctx, err.Error())
		return
	}
	isEnable, err := h.useCase.IsFeatureEnable(featureName, investorId)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[bool]{
			Data: isEnable,
		},
	)
}
