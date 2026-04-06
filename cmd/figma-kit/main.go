package main

import (
	"os"

	"github.com/amine/figma-kit/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
