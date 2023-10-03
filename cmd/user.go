package cmd

import "github.com/spf13/cobra"

func userCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "OreCast user command",
		Long: `OreCast user command
                Complete documentation is available at https://orecast.com/documentation/`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
	return cmd
}
