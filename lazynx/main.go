package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lazyengs/lazynx/internal/program"
	"github.com/lazyengs/lazynx/internal/utils"
)

func main() {
	// Setup file logging
	logger, err := utils.SetupFileLogger(utils.GetDefaultLogFile(), true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up logger: %v", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting LazyNX")

	p := program.Create()

	ctx := context.Background()
	go utils.StartNxlsclient(ctx, p, logger)

	logger.Info("Starting Bubble Tea program")
	if _, err := p.Run(); err != nil {
		logger.Errorw("Error starting program", "error", err)
		fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
		os.Exit(1)
	}

	logger.Info("LazyNX shutting down")
}
