package main

import (
	"github.com/jacobweinstock/vvalidator/cmd"
	"os"
	"os/signal"
	"syscall"
	"context"
)

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)

	defer func() {
		signal.Stop(signals)
		cancel()
	}()

	defer func() {
		signal.Stop(signals)
	}()

	go func() {
		// TODO propagate ctx through Execute
		if err := cmd.Execute(); err != nil {
			exitCode = 1
		}
	}()

	<-signals
}
