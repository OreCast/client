package cmd

import (
	"fmt"
	"os"

	oreConfig "github.com/OreCast/common/config"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile string
	verbose int

	rootCmd = &cobra.Command{
		Use:   "orecast",
		Short: "orecast command line client",
		Long: `orecast command line client
	Complete documentation is available at https://orecast.com/documentation/`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// orecast configuration
var _oreConfig *oreConfig.OreCastConfig

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.orecast.yaml)")
	rootCmd.PersistentFlags().IntVar(&verbose, "verbose", 0, "verbosity level)")

	rootCmd.AddCommand(metaCommand())
	rootCmd.AddCommand(dbsCommand())
	rootCmd.AddCommand(siteCommand())
	rootCmd.AddCommand(authCommand())
	rootCmd.AddCommand(userCommand())
}

func initConfig() {
	config, err := oreConfig.ParseConfig(cfgFile)
	if err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
	_oreConfig = &config
}
