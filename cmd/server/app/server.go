package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 1 * time.Second
)

func (app *Application) ServeHTTP() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.HttpPort),
		Handler:      app.routes(),
		ErrorLog:     log.New(os.Stderr, "", 0),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	app.Tasks.AddShutdownTask(
		func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, defaultShutdownPeriod)
			defer cancel()
			return srv.Shutdown(ctx)
		},
	)

	app.Logger.Info(fmt.Sprintf("starting server on %s", srv.Addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	app.Logger.Info(fmt.Sprintf("stopped server on %s", srv.Addr))

	return nil
}
