package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/shopspring/decimal"

	"financing-offer/cmd/server/app"
	"financing-offer/internal/version"
	"financing-offer/pkg/shutdown"
)

func main() {
	decimal.MarshalJSONWithoutQuotes = true
	showVersion := flag.Bool("version", false, "display version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	tasks, _ := shutdown.NewShutdownTasks(logger)
	defer func() {
		tasks.Wait(recover())
	}()
	err := app.Run(logger, tasks)
	if err != nil {
		trace := debug.Stack()
		logger.Error("cannot start application", slog.String("error", err.Error()), slog.String("stack", string(trace)))
		os.Exit(1)
	}
}
