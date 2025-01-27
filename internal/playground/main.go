package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/lazyengs/pkg/nxlsclient"
	"go.uber.org/zap"
)

func main() {
	_logger, _ := zap.NewDevelopment()
	logger := _logger.Sugar()
	defer logger.Sync()

	currentNxWorkspacePath, err := filepath.Abs("../..")
	if err != nil {
		logger.Fatal(err)
	}

	client := nxlsclient.NewClient(currentNxWorkspacePath, true)
	ctx, cancel := context.WithCancel(context.Background())

	// Ctrl+c like signal detection
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	ch := make(chan *nxlsclient.InitializeCommandResult)
	go client.Start(ctx, ch)

	go func() {
		<-signalChan
		logger.Infow("Received interrupt signal")
		client.Stop()
		cancel()
		signal.Stop(signalChan)
	}()

	init, ok := <-ch
	if ok {
		logger.Infow("Initialize command result", "init", init)
	}

	<-ctx.Done()
}
