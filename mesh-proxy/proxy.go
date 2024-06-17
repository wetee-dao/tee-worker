package main

import (
	"fmt"
	"os"

	"wetee.app/worker/util"
)

var version = "0.0.1"

func main() {
	if err := start(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func start() error {
	util.LogWithRed("mesh-proxy %s", version)
	return nil
}
