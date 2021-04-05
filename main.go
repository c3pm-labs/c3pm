package main

import (
	"github.com/c3pm-labs/c3pm/cmd"
)

const version = "0.0.1"
const description = `A package manager for C++.`

func main() {
	cmd.RootCmd.Version = version
	cmd.RootCmd.InitDefaultVersionFlag()
	cmd.RootCmd.Execute()
}
