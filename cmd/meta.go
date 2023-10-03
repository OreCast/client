package cmd

import "github.com/spf13/cobra"

func metaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "meta",
		Short: "OreCast meta command",
		Long: `OreCast meta command
	Complete documentation is available at https://orecast.com/documentation/`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	return cmd
}
