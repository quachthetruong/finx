package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/submissionsheet"
	"financing-offer/internal/handler"
)

type SubmissionSheetHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	UseCase submissionsheet.UseCase
}

func NewSubmissionSheetHandler(
	baseHandler handler.BaseHandler,
	logger *slog.Logger,
	useCase submissionsheet.UseCase,
) *SubmissionSheetHandler {
	return &SubmissionSheetHandler{
		BaseHandler: baseHandler,
		logger:      logger,
		UseCase:     useCase,
	}
}

func (h *SubmissionSheetHandler) Upsert(ctx *gin.Context) {
	req := entity.SubmissionSheetShorten{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.RenderParseBodyError(ctx)
		return
	}
	creator := h.UserSubOrEmpty(ctx)
	req.Metadata.Creator = creator
	res, err := h.UseCase.Upsert(ctx, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[entity.SubmissionSheet]{Data: res})
}

func (h *SubmissionSheetHandler) AdminApproveSubmissionSheet(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, "submissionId invalid")
		return
	}
	err = h.UseCase.AdminApproveSubmission(ctx, id)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[string]{Data: "success"})
}

func (h *SubmissionSheetHandler) AdminRejectSubmissionSheet(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.RenderBadRequest(ctx, "submissionId invalid")
		return
	}
	err = h.UseCase.AdminRejectSubmission(ctx, id)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, handler.BaseResponse[string]{Data: "success"})
}
