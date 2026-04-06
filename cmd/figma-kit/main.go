package main

import (
	"os"

	"github.com/dop-amine/figma-kit/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
