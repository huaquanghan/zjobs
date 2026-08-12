package main

import (
	"os"

	"zjobs/cmd/jdctl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
