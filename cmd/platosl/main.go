package main

import (
	"os"

	"github.com/platoorg/plato-sl-cli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
