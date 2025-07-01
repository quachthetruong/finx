package apperrors

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"financing-offer/cmd/server/request"
	"financing-offer/internal/appcontext"
	"financing-offer/internal/apperrors/repository"
)

type Service interface {
	// NotifyError notify error
	NotifyError(ctx context.Context, err error) error
	// Go run function in goroutine and notify error if it's not nil
	Go(ctx context.Context, f func() error)
}

type service struct {
	notifyWebhookRepository repository.NotifyWebhookRepository
	logger                  *slog.Logger
}

func (u *service) logError(ctx context.Context, err error) {
	u.logger.Error(
		"error happened", slog.String("error", err.Error()),
		slog.String("user_name", appcontext.ContextGetUserName(ctx)),
	)
}

func (u *service) Go(ctx context.Context, f func() error) {
	go func() {
		if err := f(); err != nil {
			u.logError(ctx, err)
			if err := u.NotifyError(ctx, err); err != nil {
				u.logError(ctx, err)
			}
		}
	}()
}

type NotifyFormat struct {
	UserName  string `json:"userName"`
	RequestId string `json:"requestId"`
	Trace     string `json:"trace"`
	Error     string `json:"error"`
}

func (u *service) NotifyError(ctx context.Context, err error) error {
	requestId := ""
	if ctx.Value(request.RequestIdKey) != nil {
		requestId = ctx.Value(request.RequestIdKey).(string)
	}
	notify, _ := json.Marshal(
		NotifyFormat{
			UserName:  appcontext.ContextGetUserName(ctx),
			RequestId: requestId,
			Trace:     "",
			Error:     err.Error(),
		},
	)
	if err := u.notifyWebhookRepository.Send("financing-offer-error", string(notify)); err != nil {
		return fmt.Errorf("notify error: %w", err)
	}
	return nil
}

func NewService(notifyWebhookRepository repository.NotifyWebhookRepository, logger *slog.Logger) Service {
	return &service{
		logger:                  logger,
		notifyWebhookRepository: notifyWebhookRepository,
	}
}
