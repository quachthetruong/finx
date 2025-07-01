package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/samber/do"

	"financing-offer/internal/config"
	loanOfferScheduler "financing-offer/internal/core/loanoffer/transport/scheduler"
	loanRequestScheduler "financing-offer/internal/core/loanpackagerequest/transport/scheduler"
)

var _ cron.Logger = (*logConverter)(nil)

type logConverter struct {
	base *slog.Logger
}

func (l logConverter) Info(msg string, keysAndValues ...interface{}) {
	l.base.Info(msg, keysAndValues...)
}

func (l logConverter) Error(err error, msg string, keysAndValues ...interface{}) {
	l.base.Error(fmt.Sprintf("%s, error: %s", msg, err.Error()), keysAndValues...)
}

func (app *Application) StartScheduler() error {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		return err
	}
	clog := logConverter{base: app.Logger}
	c := cron.New(
		cron.WithParser(cron.NewParser(cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
		cron.WithLocation(loc),
		cron.WithLogger(clog),
		cron.WithChain(cron.Recover(clog)),
	)
	if err := register(app.Injector, c); err != nil {
		return err
	}
	app.Tasks.AddShutdownTask(
		func(_ context.Context) error {
			cronCtx := c.Stop()
			<-cronCtx.Done()
			return nil
		},
	)
	c.Start()
	return nil
}

func register(injector *do.Injector, c *cron.Cron) error {
	cfg := do.MustInvoke[config.AppConfig](injector)
	loanOfferHandler := do.MustInvoke[*loanOfferScheduler.LoanOfferScheduler](injector)
	loanRequestHandler := do.MustInvoke[*loanRequestScheduler.LoanRequestScheduler](injector)
	if _, err := c.AddFunc(cfg.Cron.ExpireLoanOffers, loanOfferHandler.ExpireLoanOffers); err != nil {
		return err
	}
	if _, err := c.AddFunc(cfg.Cron.DeclineLoanRequests, loanRequestHandler.DeclineLoanRequests); err != nil {
		return err
	}
	return nil
}
