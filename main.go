package main

import (
	"github.com/c3pm-labs/c3pm/cmd"
)

var version = "dev"

func main() {
	cmd.RootCmd.Version = version
	cmd.RootCmd.InitDefaultVersionFlag()
	_ = cmd.RootCmd.Execute()
}
