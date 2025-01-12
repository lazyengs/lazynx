package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/lazyengs/pkg/nxlsclient"
	"go.uber.org/zap"
)

func main() {
	_logger, _ := zap.NewDevelopment()
	logger := _logger.Sugar()

	currentNxWorkspacePath, err := filepath.Abs("../..")
	if err != nil {
		logger.Fatal(err)
	}

	client := nxlsclient.NewClient(currentNxWorkspacePath, true)
	ctx, cancel := context.WithCancel(context.Background())

	// Ctrl+c like signal detection
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		client.Logger.Infow("Received interrupt signal")
		// Stops intercepting the signal
		signal.Stop(signalChan)
		cancel()
	}()

	client.Start(ctx)
}
