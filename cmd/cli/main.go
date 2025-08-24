package main

import (
	"github.io/uberate/hcli/cmd/cli/cmds"
	"os"
)

func main() {
	if err := cmds.RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
