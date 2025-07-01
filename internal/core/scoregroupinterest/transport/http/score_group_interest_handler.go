package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"financing-offer/internal/core/entity"
	scoreGroup "financing-offer/internal/core/scoregroupinterest"
	"financing-offer/internal/handler"
)

type ScoreGroupInterestHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	useCase scoreGroup.UseCase
}

const idInvalid = "id invalid"

func NewScoreGroupInterestHandler(baseHandler handler.BaseHandler, logger *slog.Logger, userUseCase scoreGroup.UseCase) *ScoreGroupInterestHandler {
	return &ScoreGroupInterestHandler{
		BaseHandler: baseHandler,
		logger:      logger,
		useCase:     userUseCase,
	}
}

func (h *ScoreGroupInterestHandler) Create(c *gin.Context) {
	u := entity.ScoreGroupInterest{}
	if err := c.ShouldBindJSON(&u); err != nil {
		h.logger.Error("create score group interest", slog.String("error", err.Error()))
		h.RenderBadRequest(c, "invalid payload")
		return
	}
	created, err := h.useCase.Create(c, u)
	if err != nil {
		h.RenderError(c, err)
		return
	}
	c.JSON(
		201, handler.BaseResponse[entity.ScoreGroupInterest]{
			Data: created,
		},
	)
}

func (h *ScoreGroupInterestHandler) Update(c *gin.Context) {
	id, err := h.ParamsInt(c)
	errorMessage := "update score group interest"
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(c, idInvalid)
		return
	}
	role, err := h.useCase.GetById(c, id)
	if err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderNotFound(c, fmt.Sprintf("group role id: %v not found", id))
		return
	}
	if err := c.ShouldBindJSON(&role); err != nil {
		h.logger.Error(errorMessage, slog.String("error", err.Error()))
		h.RenderBadRequest(c, "invalid payload")
		return
	}

	updated, err := h.useCase.Update(c, role)
	if err != nil {
		h.RenderError(c, err)
		return
	}
	c.JSON(
		200, handler.BaseResponse[entity.ScoreGroupInterest]{
			Data: updated,
		},
	)
}

func (h *ScoreGroupInterestHandler) GetAll(c *gin.Context) {
	roles, err := h.useCase.GetAll(c)
	if err != nil {
		h.RenderError(c, err)
		return
	}
	sort.SliceStable(
		roles, func(i, j int) bool {
			return roles[i].Id < roles[j].Id
		},
	)
	c.JSON(
		200, handler.BaseResponse[[]entity.ScoreGroupInterest]{
			Data: roles,
		},
	)
}

func (h *ScoreGroupInterestHandler) Delete(c *gin.Context) {
	id, err := h.ParamsInt(c)
	if err != nil {
		h.logger.Error("delete score group interest", slog.String("error", err.Error()))
		h.RenderBadRequest(c, idInvalid)
		return
	}
	isDeleted, err := h.useCase.Delete(c, id)
	if err != nil {
		h.RenderError(c, err)
		return
	}
	c.JSON(
		200, handler.BaseResponse[bool]{
			Data: isDeleted,
		},
	)
}

func (h *ScoreGroupInterestHandler) GetById(ctx *gin.Context) {
	id, err := h.ParamsInt(ctx)
	if err != nil {
		h.logger.Error("get score group interest", slog.String("error", err.Error()))
		h.RenderBadRequest(ctx, idInvalid)
		return
	}
	res, err := h.useCase.GetById(ctx, id)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		http.StatusOK, handler.BaseResponse[entity.ScoreGroupInterest]{
			Data: res,
		},
	)
}
