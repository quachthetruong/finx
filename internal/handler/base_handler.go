package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"financing-offer/cmd/server/request"
	"financing-offer/internal/appcontext"
	"financing-offer/internal/apperrors"
	"financing-offer/internal/core"
	"financing-offer/internal/core/entity"
	"financing-offer/pkg/number"
)

const (
	pageSizeMax = 1000
)

type BaseHandler struct {
	Logger       *slog.Logger
	ErrorUseCase apperrors.Service
}

func NewBaseHandler(logger *slog.Logger, errorUseCase apperrors.Service) BaseHandler {
	return BaseHandler{Logger: logger, ErrorUseCase: errorUseCase}
}

func (h *BaseHandler) ParamsInt(c *gin.Context) (int64, error) {
	idParam := c.Param("id")
	return strconv.ParseInt(idParam, 10, 64)
}

func (h *BaseHandler) ParamNotEmpty(c *gin.Context, key string) (string, error) {
	value := c.Param(key)
	if value == "" {
		return "", apperrors.ErrParamInvalid(key)
	}
	return value, nil
}

func (h *BaseHandler) InvestorId(c *gin.Context) (string, error) {
	customerInfo := appcontext.ContextGetCustomerInfo(c)
	if customerInfo != nil && customerInfo.InvestorId != "" {
		return customerInfo.InvestorId, nil
	}
	return "", apperrors.ErrInvalidInvestorId
}

func (h *BaseHandler) Investor(c *gin.Context) (entity.Investor, error) {
	customerInfo := appcontext.ContextGetCustomerInfo(c)
	if customerInfo != nil {
		return entity.Investor{
			InvestorId:  customerInfo.InvestorId,
			CustodyCode: customerInfo.CustodyCode,
		}, nil
	}
	return entity.Investor{}, apperrors.ErrInvalidInvestorId
}

func (h *BaseHandler) UserSubOrEmpty(c *gin.Context) string {
	customerInfo := appcontext.ContextGetCustomerInfo(c)
	if customerInfo != nil {
		return customerInfo.Sub
	}
	return ""
}

func (h *BaseHandler) RenderUnauthenticated(c *gin.Context, messages ...string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": strings.Join(messages, ",")})
}

func (h *BaseHandler) RenderBadRequest(c *gin.Context, messages ...string) {
	c.JSON(
		http.StatusBadRequest, gin.H{
			"error": strings.Join(messages, ","),
			"code":  http.StatusBadRequest,
		},
	)
}

func (h *BaseHandler) RenderIdInvalid(c *gin.Context) {
	h.RenderBadRequest(c, "id invalid")
}

func (h *BaseHandler) RenderParseBodyError(c *gin.Context) {
	h.RenderBadRequest(c, "parse body")
}

func (h *BaseHandler) RenderNotFound(c *gin.Context, messages ...string) {
	c.JSON(http.StatusNotFound, gin.H{"error": strings.Join(messages, ",")})
}

func (h *BaseHandler) RenderError(c *gin.Context, err error) {
	var appErr apperrors.AppError
	requestId := ""
	requestIdVal, _ := c.Get(request.RequestIdKey)
	if requestIdString, ok := requestIdVal.(string); ok {
		requestId = requestIdString
	}
	h.Logger.Error(
		"error happened while handling request", slog.String("error", err.Error()),
		slog.String("user_name", appcontext.ContextGetUserName(c)), slog.String("request_id", requestId),
	)
	if !apperrors.IsNotFoundError(err) {
		h.ReportError(c, err)
	}
	if ok := errors.As(err, &appErr); ok {
		code := number.GetFirstThreeDigits(appErr.Code)
		if code >= http.StatusInternalServerError {
			c.JSON(code, gin.H{"error": "an error happened, please try again later"})
			return
		}
		c.JSON(
			code, gin.H{
				"error": appErr.Message,
				"code":  appErr.Code,
			},
		)
		return
	}
	// handle non-AppError errors
	if apperrors.IsConstraintViolationError(err) {
		c.JSON(http.StatusConflict, gin.H{"error": "constraint violation"})
		return
	}
	if apperrors.IsNotFoundError(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found resources"})
		return
	}
	if apperrors.IsObjectNotInPrerequisiteStateError(err) {
		c.JSON(
			http.StatusLocked,
			gin.H{"error": "the requested resource is currently unavailable, please try again later"},
		)
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "an error happened, please try again later"})
}

func (h *BaseHandler) ReportError(ctx context.Context, err error) {
	if notifyErr := h.ErrorUseCase.NotifyError(ctx, err); notifyErr != nil {
		h.Logger.Error(
			"error happened while notifying error", slog.String("error", err.Error()),
			slog.String("notify_error", notifyErr.Error()),
		)
	}
}

func (h *BaseHandler) ParseQueryWithPagination(ctx *gin.Context, paging *core.Paging, req any) error {
	if err := ctx.ShouldBindQuery(req); err != nil {
		return err
	}
	if err := h.parsePagination(ctx, paging); err != nil {
		return err
	}
	return nil
}

func (h *BaseHandler) parsePagination(ctx *gin.Context, paging *core.Paging) error {
	paging.Number = 1

	if sizeStr := ctx.Query("page[size]"); sizeStr != "" {
		size, err := strconv.Atoi(sizeStr)
		if err != nil || size == 0 || uint(size) > pageSizeMax {
			return apperrors.ErrParamInvalid("page[size]")
		}
		paging.Size = uint(size)
	}
	if numberStr := ctx.Query("page[number]"); numberStr != "" {
		pageNumber, err := strconv.Atoi(numberStr)
		if err != nil || pageNumber < 0 {
			return apperrors.ErrParamInvalid("page[number]")
		}
		if pageNumber > 0 {
			paging.Number = uint(pageNumber)
		}
	}
	paging.Sort = parseSort(ctx, paging.Sort)
	return nil
}

func parseSort(ctx *gin.Context, original core.Orders) core.Orders {
	orders := core.Orders{}
	if sortQuery := ctx.Query("sort"); sortQuery != "" {
		sort := strings.Split(sortQuery, ",")
		for _, str := range sort {
			if strings.HasPrefix(str, "-") {
				if len(str) == 1 {
					continue
				}
				orders.Add(core.Order{Direction: core.DirectionDesc, ColumnName: str[1:]})
			} else {
				orders.Add(core.Order{Direction: core.DirectionAsc, ColumnName: str})
			}
		}
	}
	if len(original) > 0 {
		orders.Add(original...)
	}
	return orders
}
