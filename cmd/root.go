// Package cmd hosts the configuration and handling of the command line interface of C3PM.
package cmd

import "github.com/alecthomas/kong"

//CLI is the root configuration of C3PM's command line interface.
var CLI struct {
	Version kong.VersionFlag `short:"v" help:"outputs the version number"`
	Add     AddCmd           `kong:"cmd,help='Add a new dependency'"`
	Init    InitCmd          `kong:"cmd,help='Init a c3pm project'"`
	Logout  LogoutCmd        `kong:"cmd,help='Logout from the api'"`
	Login   LoginCmd         `kong:"cmd,help='Login to the api'"`
	Build   BuildCmd         `kong:"cmd,help='Build a c3pm project'"`
	Publish PublishCmd       `kong:"cmd,help='Publish a c3pm project'"`
}
