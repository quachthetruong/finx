package database

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/lib/pq"

	"financing-offer/internal/apperrors"
	"financing-offer/internal/config"
	"financing-offer/pkg/shutdown"
)

type DbListener struct {
	logger        *slog.Logger
	wg            sync.WaitGroup
	dsn           string
	tasks         *shutdown.Tasks
	errorReporter apperrors.Service
}

func (l *DbListener) Listen(ctx context.Context, channel string, callback func(ctx context.Context, data string) error) error {
	listener := pq.NewListener(l.dsn, 10*time.Second, time.Minute, nil)
	if err := listener.Listen(channel); err != nil {
		return err
	}
	l.tasks.AddShutdownTask(
		func(_ context.Context) error {
			return listener.Close()
		},
	)
	l.logger.Info("start listening to db notify", slog.String("channel", channel))
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case n := <-listener.Notify:
				l.logger.Info("received db event", slog.String("event", n.Extra))
				if n == nil {
					continue
				}
				l.wg.Add(1)
				go func() {
					defer l.wg.Done()
					defer func() {
						if r := recover(); r != nil {
							l.logger.Error("panic recovered from db listener", slog.Any("panic", r))
							_ = l.errorReporter.NotifyError(ctx, fmt.Errorf("panic recovered from db listener"))
						}
					}()
					if err := callback(context.Background(), n.Extra); err != nil {
						l.logger.Error("error while handling db event", slog.String("error", err.Error()))
						_ = l.errorReporter.NotifyError(ctx, err)
					}
				}()
			}
		}
	}()
	return nil
}

func NewListener(logger *slog.Logger, cfg config.DbConfig, tasks *shutdown.Tasks, errorReporter apperrors.Service) *DbListener {
	completeDsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?binary_parameters=yes", cfg.User, cfg.Password,
		cfg.Host, cfg.Port, cfg.DbName,
	)
	if !cfg.EnableSsl {
		completeDsn += "&sslmode=disable"
	}
	listener := &DbListener{
		errorReporter: errorReporter,
		logger:        logger,
		dsn:           completeDsn,
		tasks:         tasks,
	}
	tasks.AddShutdownTask(
		func(_ context.Context) error {
			listener.wg.Wait()
			return nil
		},
	)
	return listener
}
