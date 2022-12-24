package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "gotcha",
	Short: "Pull Request Quality Reviewer",
}

// Starts the root-cli.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
