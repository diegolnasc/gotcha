package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version configuration.
type GotchaVersion struct {
	Version string
}

// Returns the command cli
func (v *GotchaVersion) Init() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the current version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Gotcha -> " + v.Version)
		},
	}
}
