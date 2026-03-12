package main

import (
	"fmt"
	"os"

	"github.com/Epsilondelta-ai/moai-adk-codex/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:], os.Stdout, os.Stderr, os.Getwd); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
