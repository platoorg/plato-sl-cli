package main

import (
	"os"

	"platosl.org/cmd/platosl/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
