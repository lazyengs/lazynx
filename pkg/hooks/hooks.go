package hooks

import (
	"errors"
	"log"
	"os"
)

func EnsureNxProject() {
	if _, err := os.Stat("nx.json"); errors.Is(err, os.ErrNotExist) {
		log.Println("Error: must be run inside an Nx project")

		os.Exit(1)
	}
}
