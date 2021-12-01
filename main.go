package main

import (
	"github.com/c3pm-labs/c3pm/cmd"
)

var version = "1.0.0"

func main() {
	cmd.RootCmd.Version = version
	cmd.RootCmd.InitDefaultVersionFlag()
	_ = cmd.RootCmd.Execute()
}
