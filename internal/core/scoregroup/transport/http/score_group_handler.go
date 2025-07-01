package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/scoregroup"
	"financing-offer/internal/handler"
)

type ScoreGroupHandler struct {
	handler.BaseHandler
	Logger            *slog.Logger
	scoreGroupUseCase scoregroup.UseCase
}

func NewScoreGroupHandler(
	baseHandler handler.BaseHandler,
	logger *slog.Logger,
	scoreGroupUseCase scoregroup.UseCase,
) *ScoreGroupHandler {
	return &ScoreGroupHandler{
		BaseHandler:       baseHandler,
		Logger:            logger,
		scoreGroupUseCase: scoreGroupUseCase,
	}
}

func (h *ScoreGroupHandler) GetAll(ctx *gin.Context) {
	scoreGroups, err := h.scoreGroupUseCase.GetAll(ctx)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.ScoreGroup]{
			Data: scoreGroups,
		},
	)
}

func (h *ScoreGroupHandler) Create(ctx *gin.Context) {
	req := ScoreGroupRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("create score group", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	created, err := h.scoreGroupUseCase.Create(ctx, req.toEntity(0))
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusCreated, handler.BaseResponse[entity.ScoreGroup]{
			Data: created,
		},
	)
}

func (h *ScoreGroupHandler) Update(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.Logger.Error("update score group", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	oldScoreGroup, err := h.scoreGroupUseCase.GetById(ctx, id)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	if err := ctx.ShouldBindJSON(&oldScoreGroup); err != nil {
		h.Logger.Error("update score group", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid payload")
		return
	}
	updated, err := h.scoreGroupUseCase.Update(ctx, oldScoreGroup)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.ScoreGroup]{
			Data: updated,
		},
	)
}

func (h *ScoreGroupHandler) Delete(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.Logger.Error("delete score group", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "id invalid")
		return
	}
	if err := h.scoreGroupUseCase.Delete(ctx, id); err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusNoContent, handler.BaseResponse[string]{
			Data: "ok",
		},
	)
}

func (h *ScoreGroupHandler) GetAvailablePackages(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.Logger.Error("get available package for score group", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, "invalid id")
		return
	}
	availablePackages, err := h.scoreGroupUseCase.GetAvailablePackagesById(ctx, id)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[[]entity.ScoreGroupInterest]{
			Data: availablePackages,
		},
	)
}
