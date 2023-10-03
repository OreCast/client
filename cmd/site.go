package cmd

import "github.com/spf13/cobra"

func siteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "OreCast site command",
		Long: `OreCast site command
                Complete documentation is available at https://orecast.com/documentation/`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	return cmd
}
