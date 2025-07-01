package http

import (
	"financing-offer/internal/config"
	"financing-offer/internal/handler"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type ConfigHandler struct {
	handler.BaseHandler
	logger  *slog.Logger
	useCase config.UseCase
}

func NewConfigHandler(
	baseHandler handler.BaseHandler,
	logger *slog.Logger,
	useCase config.UseCase,
) *ConfigHandler {
	return &ConfigHandler{
		BaseHandler: baseHandler,
		logger:      logger,
		useCase:     useCase,
	}
}

// GetConfiguration godoc
//
//	@Summary		Get configuration
//	@Description	Get configuration
//	@Tags			configuration,investor
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	handler.BaseResponse[map[string]string]
//	@Failure		500	{object}	handler.ErrorResponse
//	@Security		BearerAuth
//	@Router			/v1/my-configurations [get]
func (h *ConfigHandler) GetConfiguration(ctx *gin.Context) {
	res, err := h.useCase.GetConfigurations()
	if err != nil {
		h.RenderError(ctx, err)
		return
	}
	ctx.JSON(
		200, handler.BaseResponse[map[string]any]{
			Data: res,
		},
	)
}
