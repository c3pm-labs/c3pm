package main

import (
	"github.com/alecthomas/kong"
	"github.com/c3pm-labs/c3pm/cmd"
)

const version = "0.0.1"
const description = `A package manager for C++.`

func main() {
	ctx := kong.Parse(&cmd.CLI, kong.Description(description), kong.UsageOnError(), kong.Vars{
		"version": version,
	})
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
