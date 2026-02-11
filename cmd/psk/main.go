package main

import (
	"os"

	"github.com/c8ab/provenskills/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args))
}
