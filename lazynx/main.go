package main

import (
	"fmt"
	"os"

	"github.com/lazyengs/lazynx/internal/program"
)

func main() {
	p := program.Create()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
		os.Exit(1)
	}
}
