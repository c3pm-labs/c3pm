// Package cmd hosts the configuration and handling of the command line interface of C3PM.
package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ctpm",
	Short: "c3pm abstracts your build system and eases the management of your dependencies.",
	Long:  "C3PM is a next-generation package manager for C++.\nYou can use C3PM to share and use packages with other developers around the world.",
}

func init() {
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(buildCmd)
	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(logoutCmd)
	RootCmd.AddCommand(loginCmd)
	RootCmd.AddCommand(publishCmd)
	RootCmd.AddCommand(removeCmd)
	RootCmd.AddCommand(testCmd)
}
