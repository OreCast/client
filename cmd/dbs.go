package cmd

import "github.com/spf13/cobra"

func dbsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dbs",
		Short: "OreCast dbs command",
		Long: `OreCast data-bookkeeping system command
                Complete documentation is available at https://orecast.com/documentation/`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	return cmd
}
