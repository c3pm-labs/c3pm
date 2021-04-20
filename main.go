package main

import (
	"github.com/c3pm-labs/c3pm/cmd"
)

<<<<<<< HEAD
var version = "dev"
=======
var version = "0.0.1"
>>>>>>> 0d7d584b99cc442c325a8cdb836a11be87c58947

func main() {
	cmd.RootCmd.Version = version
	cmd.RootCmd.InitDefaultVersionFlag()
	_ = cmd.RootCmd.Execute()
}
